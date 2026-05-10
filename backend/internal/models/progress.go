package models

import "time"

type LessonProgress struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	LessonID  string    `json:"lesson_id"`
	Status    string    `json:"status"` // not_started, in_progress, completed
	UpdatedAt time.Time `json:"updated_at"`
}

// Runner models (for inter-service communication)

type RunRequest struct {
	SourceCode string         `json:"source_code"`
	TestCases  []RunTestCase  `json:"test_cases"`
	Mode       string         `json:"mode"` // "judge" or "run"
}

type RunTestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type RunResult struct {
	Verdict        string `json:"verdict"`
	TestsPassed    int    `json:"tests_passed"`
	TestsTotal     int    `json:"tests_total"`
	ExecTimeMs     int    `json:"exec_time_ms"`
	CompilerOutput string `json:"compiler_output"`
	FailedTest     *int   `json:"failed_test"`
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
}

// XP award models

type XPAwardRequest struct {
	UserID        string `json:"user_id"`
	SolutionID    string `json:"solution_id"`
	TaskID        string `json:"task_id"`
	Difficulty    string `json:"difficulty"`
	AttemptNumber int    `json:"attempt_number"`
}

type XPAwardResponse struct {
	XPAwarded        int           `json:"xp_awarded"`
	TotalXP          int           `json:"total_xp"`
	Level            int           `json:"level"`
	NewAchievements  []Achievement `json:"new_achievements"`
	AlreadyProcessed bool          `json:"already_processed,omitempty"`
}
