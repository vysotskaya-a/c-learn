package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	LessonID    string    `json:"lesson_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskDetail struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Difficulty  string       `json:"difficulty"`
	IsSolved    bool         `json:"is_solved"`
	Samples     []SampleTest `json:"samples"`
}

type TestCase struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	Input     string `json:"input"`
	Expected  string `json:"expected"`
	SortOrder int    `json:"sort_order"`
	IsSample  bool   `json:"is_sample"`
}

type SampleTest struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}
