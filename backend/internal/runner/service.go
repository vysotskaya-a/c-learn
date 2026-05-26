package runner

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/c-learn/internal/models"
)

const (
	VerdictOK                = "ok"
	VerdictCompilationError  = "compilation_error"
	VerdictWrongAnswer       = "wrong_answer"
	VerdictTimeLimitExceeded = "time_limit_exceeded"
	VerdictRuntimeError      = "runtime_error"
)

type containerRunner interface {
	CreateContainer(ctx context.Context) (string, error)
	RemoveContainer(ctx context.Context, containerID string) error
	CopySourceToContainer(ctx context.Context, containerID, sourceCode string) error
	Compile(ctx context.Context, containerID string) (*ExecResult, error)
	RunTest(ctx context.Context, containerID, input string) (*ExecResult, error)
}

type Service struct {
	docker containerRunner
}

func NewService(docker *DockerRunner) *Service {
	return &Service{docker: docker}
}

func (s *Service) Run(ctx context.Context, req models.RunRequest) (*models.RunResult, error) {
	start := time.Now()

	containerID, err := s.docker.CreateContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}
	defer func() {
		if rmErr := s.docker.RemoveContainer(context.Background(), containerID); rmErr != nil {
			log.Printf("remove container %s: %v", containerID[:12], rmErr)
		}
	}()

	// Copy source code into container
	if err := s.docker.CopySourceToContainer(ctx, containerID, req.SourceCode); err != nil {
		return nil, fmt.Errorf("copy source: %w", err)
	}

	// Compile
	compileResult, err := s.docker.Compile(ctx, containerID)
	if err != nil {
		return &models.RunResult{
			Verdict:        VerdictCompilationError,
			CompilerOutput: fmt.Sprintf("compilation error: %v", err),
			TestsTotal:     len(req.TestCases),
		}, nil
	}
	if compileResult.ExitCode != 0 {
		return &models.RunResult{
			Verdict:        VerdictCompilationError,
			CompilerOutput: compileResult.Stderr,
			TestsTotal:     len(req.TestCases),
		}, nil
	}

	// Mode "run": single execution with stdin, no comparison
	if req.Mode == "run" {
		stdin := ""
		if len(req.TestCases) > 0 {
			stdin = req.TestCases[0].Input
		}
		runResult, err := s.docker.RunTest(ctx, containerID, stdin)
		if err != nil {
			if strings.Contains(err.Error(), "time_limit_exceeded") {
				return &models.RunResult{
					Verdict:    VerdictTimeLimitExceeded,
					ExecTimeMs: int(time.Since(start).Milliseconds()),
				}, nil
			}
			return &models.RunResult{
				Verdict: VerdictRuntimeError,
				Stderr:  fmt.Sprintf("%v", err),
			}, nil
		}
		return &models.RunResult{
			Verdict:    VerdictOK,
			Stdout:     runResult.Stdout,
			Stderr:     runResult.Stderr,
			ExecTimeMs: int(time.Since(start).Milliseconds()),
		}, nil
	}

	// Mode "judge": run all test cases
	for i, tc := range req.TestCases {
		runResult, err := s.docker.RunTest(ctx, containerID, tc.Input)
		if err != nil {
			if strings.Contains(err.Error(), "time_limit_exceeded") {
				return &models.RunResult{
					Verdict:     VerdictTimeLimitExceeded,
					TestsPassed: i,
					TestsTotal:  len(req.TestCases),
					ExecTimeMs:  int(time.Since(start).Milliseconds()),
				}, nil
			}
			return &models.RunResult{
				Verdict:     VerdictRuntimeError,
				TestsPassed: i,
				TestsTotal:  len(req.TestCases),
				ExecTimeMs:  int(time.Since(start).Milliseconds()),
			}, nil
		}

		if runResult.ExitCode != 0 {
			return &models.RunResult{
				Verdict:     VerdictRuntimeError,
				TestsPassed: i,
				TestsTotal:  len(req.TestCases),
				ExecTimeMs:  int(time.Since(start).Milliseconds()),
			}, nil
		}

		if NormalizeOutput(runResult.Stdout) != NormalizeOutput(tc.Expected) {
			failedTest := i + 1
			return &models.RunResult{
				Verdict:     VerdictWrongAnswer,
				TestsPassed: i,
				TestsTotal:  len(req.TestCases),
				FailedTest:  &failedTest,
				ExecTimeMs:  int(time.Since(start).Milliseconds()),
			}, nil
		}
	}

	return &models.RunResult{
		Verdict:     VerdictOK,
		TestsPassed: len(req.TestCases),
		TestsTotal:  len(req.TestCases),
		ExecTimeMs:  int(time.Since(start).Milliseconds()),
	}, nil
}
