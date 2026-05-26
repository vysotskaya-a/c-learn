package runner

import (
	"context"
	"fmt"
	"testing"

	"github.com/c-learn/internal/models"
)

type mockContainerRunner struct {
	compileExit int
	compileErr  string
	runs        []ExecResult
	runErr      error
}

func (m *mockContainerRunner) CreateContainer(context.Context) (string, error) {
	return "container-1", nil
}

func (m *mockContainerRunner) RemoveContainer(context.Context, string) error {
	return nil
}

func (m *mockContainerRunner) CopySourceToContainer(context.Context, string, string) error {
	return nil
}

func (m *mockContainerRunner) Compile(context.Context, string) (*ExecResult, error) {
	return &ExecResult{
		Stdout:   "",
		Stderr:   m.compileErr,
		ExitCode: m.compileExit,
	}, nil
}

func (m *mockContainerRunner) RunTest(_ context.Context, _ string, _ string) (*ExecResult, error) {
	if m.runErr != nil {
		return nil, m.runErr
	}
	if len(m.runs) == 0 {
		return &ExecResult{Stdout: "", ExitCode: 0}, nil
	}
	r := m.runs[0]
	m.runs = m.runs[1:]
	return &r, nil
}

func TestService_Run_ModeRun(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 0,
		runs:        []ExecResult{{Stdout: "hello\n", Stderr: "", ExitCode: 0}},
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "#include <stdio.h>",
		TestCases:  []models.RunTestCase{{Input: "ignored"}},
		Mode:       "run",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictOK || result.Stdout != "hello\n" {
		t.Fatalf("result = %+v", result)
	}
}

func TestService_Run_CompilationError(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 1,
		compileErr:  "syntax error",
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "invalid",
		TestCases:  []models.RunTestCase{{Input: "1", Expected: "1"}},
		Mode:       "judge",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictCompilationError || result.TestsTotal != 1 {
		t.Fatalf("result = %+v", result)
	}
}

func TestService_Run_JudgeOK(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 0,
		runs: []ExecResult{
			{Stdout: "3\n", ExitCode: 0},
			{Stdout: "5\n", ExitCode: 0},
		},
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "code",
		TestCases: []models.RunTestCase{
			{Input: "", Expected: "3"},
			{Input: "", Expected: "5"},
		},
		Mode: "judge",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictOK || result.TestsPassed != 2 || result.TestsTotal != 2 {
		t.Fatalf("result = %+v", result)
	}
}

func TestService_Run_WrongAnswer(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 0,
		runs:        []ExecResult{{Stdout: "99\n", ExitCode: 0}},
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "code",
		TestCases:  []models.RunTestCase{{Input: "", Expected: "42"}},
		Mode:       "judge",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictWrongAnswer || result.TestsPassed != 0 {
		t.Fatalf("result = %+v", result)
	}
	if result.FailedTest == nil || *result.FailedTest != 1 {
		t.Fatalf("FailedTest = %v, want 1", result.FailedTest)
	}
}

func TestService_Run_TimeLimitExceeded(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 0,
		runErr:      fmt.Errorf("time_limit_exceeded"),
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "code",
		TestCases:  []models.RunTestCase{{Input: "", Expected: "1"}},
		Mode:       "judge",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictTimeLimitExceeded {
		t.Fatalf("result = %+v", result)
	}
}

func TestService_Run_RuntimeError(t *testing.T) {
	mock := &mockContainerRunner{
		compileExit: 0,
		runs:        []ExecResult{{Stdout: "", ExitCode: 139}},
	}
	svc := &Service{docker: mock}

	result, err := svc.Run(context.Background(), models.RunRequest{
		SourceCode: "code",
		TestCases:  []models.RunTestCase{{Input: "", Expected: "1"}},
		Mode:       "judge",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Verdict != VerdictRuntimeError {
		t.Fatalf("result = %+v", result)
	}
}
