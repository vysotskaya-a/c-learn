package runner

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	defaultImage   = "gcc:13"
	memoryLimit    = 64 * 1024 * 1024 // 64 MB
	pidsLimit      = 50
	compileTimeout = 10 * time.Second
	runTimeout     = 2 * time.Second
	maxOutputSize  = 1 * 1024 * 1024 // 1 MB
	workDir        = "/tmp/work"
)

type DockerRunner struct {
	cli       *client.Client
	imageName string
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// NewDockerRunner creates a new runner and ensures the sandbox image is available.
// imageName can be empty — defaults to "gcc:13".
func NewDockerRunner(imageName string) (*DockerRunner, error) {
	if imageName == "" {
		imageName = defaultImage
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	runner := &DockerRunner{cli: cli, imageName: imageName}

	// Ensure sandbox image is available — pull if missing
	if err := runner.ensureImage(context.Background()); err != nil {
		return nil, fmt.Errorf("ensure sandbox image: %w", err)
	}

	return runner, nil
}

// ensureImage checks if the image exists locally; if not, pulls it.
func (d *DockerRunner) ensureImage(ctx context.Context) error {
	filterArgs := filters.NewArgs()
	filterArgs.Add("reference", d.imageName)

	images, err := d.cli.ImageList(ctx, types.ImageListOptions{Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("list images: %w", err)
	}

	if len(images) > 0 {
		log.Printf("[runner] sandbox image %s found locally", d.imageName)
		return nil
	}

	log.Printf("[runner] sandbox image %s not found, pulling...", d.imageName)
	reader, err := d.cli.ImagePull(ctx, "docker.io/library/"+d.imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pull image %s: %w", d.imageName, err)
	}
	defer reader.Close()

	// Read the pull output to completion (required to finish the pull)
	_, _ = io.Copy(io.Discard, reader)

	log.Printf("[runner] sandbox image %s pulled successfully", d.imageName)
	return nil
}

func (d *DockerRunner) CreateContainer(ctx context.Context) (string, error) {
	pids := int64(pidsLimit)

	cfg := &container.Config{
		Image:           d.imageName,
		Cmd:             []string{"sleep", "60"},
		Tty:             false,
		NetworkDisabled: true,
		WorkingDir:      workDir,
	}

	hostCfg := &container.HostConfig{
		Resources: container.Resources{
			Memory:     memoryLimit,
			MemorySwap: memoryLimit,
			CPUPeriod:  100000,
			CPUQuota:   200000,
			PidsLimit:  &pids,
		},
		ReadonlyRootfs: false,
		SecurityOpt:    []string{"no-new-privileges"},
		NetworkMode:    "none",
	}

	resp, err := d.cli.ContainerCreate(ctx, cfg, hostCfg, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("create container: %w", err)
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		_ = d.RemoveContainer(ctx, resp.ID)
		return "", fmt.Errorf("start container: %w", err)
	}

	// Create work directory inside the running container
	mkdirResult, err := d.Exec(ctx, resp.ID, []string{"mkdir", "-p", workDir})
	if err != nil {
		_ = d.RemoveContainer(ctx, resp.ID)
		return "", fmt.Errorf("mkdir workspace: %w", err)
	}
	if mkdirResult.ExitCode != 0 {
		_ = d.RemoveContainer(ctx, resp.ID)
		return "", fmt.Errorf("mkdir workspace failed: %s", mkdirResult.Stderr)
	}

	return resp.ID, nil
}

func (d *DockerRunner) RemoveContainer(ctx context.Context, id string) error {
	return d.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true})
}

// CopySourceToContainer writes source code into the container via exec + base64.
// More reliable than docker cp, especially on macOS Docker Desktop.
func (d *DockerRunner) CopySourceToContainer(ctx context.Context, containerID string, content string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	// Write file via shell: echo <base64> | base64 -d > solution.c
	cmd := []string{
		"sh", "-c",
		fmt.Sprintf("echo '%s' | base64 -d > %s/solution.c", encoded, workDir),
	}

	result, err := d.Exec(ctx, containerID, cmd)
	if err != nil {
		return fmt.Errorf("write source file: %w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("write source file failed (exit %d): %s", result.ExitCode, result.Stderr)
	}

	// Verify file exists
	verifyResult, err := d.Exec(ctx, containerID, []string{"test", "-f", workDir + "/solution.c"})
	if err != nil || verifyResult.ExitCode != 0 {
		return fmt.Errorf("source file verification failed: file not found in container")
	}

	return nil
}

func (d *DockerRunner) Exec(ctx context.Context, containerID string, cmd []string) (*ExecResult, error) {
	execCfg := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := d.cli.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return nil, fmt.Errorf("exec create: %w", err)
	}

	resp, err := d.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, fmt.Errorf("exec attach: %w", err)
	}
	defer resp.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	limitedReader := io.LimitReader(resp.Reader, maxOutputSize)
	_, _ = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, limitedReader)

	inspect, err := d.cli.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return &ExecResult{
			Stdout: stdoutBuf.String(),
			Stderr: stderrBuf.String(),
		}, nil
	}

	return &ExecResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: inspect.ExitCode,
	}, nil
}

func (d *DockerRunner) ExecWithStdin(ctx context.Context, containerID string, cmd []string, stdin string) (*ExecResult, error) {
	execCfg := types.ExecConfig{
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := d.cli.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return nil, fmt.Errorf("exec create: %w", err)
	}

	resp, err := d.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, fmt.Errorf("exec attach: %w", err)
	}
	defer resp.Close()

	if stdin != "" {
		_, _ = resp.Conn.Write([]byte(stdin))
	}
	resp.CloseWrite()

	var stdoutBuf, stderrBuf bytes.Buffer
	limitedReader := io.LimitReader(resp.Reader, maxOutputSize)
	_, _ = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, limitedReader)

	inspect, err := d.cli.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return &ExecResult{
			Stdout: stdoutBuf.String(),
			Stderr: stderrBuf.String(),
		}, nil
	}

	return &ExecResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: inspect.ExitCode,
	}, nil
}

func (d *DockerRunner) Compile(ctx context.Context, containerID string) (*ExecResult, error) {
	compileCtx, cancel := context.WithTimeout(ctx, compileTimeout)
	defer cancel()
	return d.Exec(compileCtx, containerID, []string{
		"gcc", "-o", workDir + "/solution", workDir + "/solution.c", "-lm", "-Wall",
	})
}

func (d *DockerRunner) RunTest(ctx context.Context, containerID, input string) (*ExecResult, error) {
	runCtx, cancel := context.WithTimeout(ctx, runTimeout)
	defer cancel()
	result, err := d.ExecWithStdin(runCtx, containerID, []string{workDir + "/solution"}, input)
	if err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("time_limit_exceeded")
		}
		return nil, err
	}
	return result, nil
}

func NormalizeOutput(s string) string {
	return strings.TrimRight(s, "\n")
}
