package models

import "time"

type Lesson struct {
	ID        string    `json:"id"`
	ModuleID  string    `json:"module_id"`
	Title     string    `json:"title"`
	TheoryMD  string    `json:"theory_md"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LessonBrief struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
	Status    string `json:"status"`
}

type LessonDetail struct {
	ID       string       `json:"id"`
	ModuleID string       `json:"module_id"`
	Title    string       `json:"title"`
	TheoryMD string       `json:"theory_md"`
	Status   string       `json:"status"`
	Tasks    []TaskDetail `json:"tasks"`
}
