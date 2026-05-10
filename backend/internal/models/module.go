package models

import "time"

type Module struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ModuleWithLessons struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	SortOrder int            `json:"sort_order"`
	Lessons   []LessonBrief  `json:"lessons"`
}
