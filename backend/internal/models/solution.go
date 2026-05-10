package models

import "time"

type Solution struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TaskID         string    `json:"task_id"`
	SourceCode     string    `json:"source_code"`
	Verdict        string    `json:"verdict"`
	CompilerOutput string    `json:"compiler_output,omitempty"`
	TestsPassed    int       `json:"tests_passed"`
	TestsTotal     int       `json:"tests_total"`
	ExecTimeMs     int       `json:"exec_time_ms,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type SubmitResponse struct {
	SolutionID           string        `json:"solution_id"`
	Verdict              string        `json:"verdict"`
	TestsPassed          int           `json:"tests_passed"`
	TestsTotal           int           `json:"tests_total"`
	ExecTimeMs           int           `json:"exec_time_ms,omitempty"`
	CompilerOutput       string        `json:"compiler_output,omitempty"`
	FailedTest           int           `json:"failed_test,omitempty"`
	XPAwarded            int           `json:"xp_awarded"`
	AchievementsUnlocked []Achievement `json:"achievements_unlocked"`
}

type RunResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	ExecTimeMs int    `json:"exec_time_ms"`
}
