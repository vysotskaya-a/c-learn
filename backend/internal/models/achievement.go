package models

import "time"

type Achievement struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
}

type UserAchievement struct {
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url,omitempty"`
	AwardedAt   time.Time `json:"awarded_at"`
}

type Profile struct {
	UserID       string            `json:"user_id"`
	TotalXP      int               `json:"total_xp"`
	Level        int               `json:"level"`
	SolvedCount  int               `json:"solved_count"`
	Achievements []UserAchievement `json:"achievements"`
}

type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	TotalXP     int    `json:"total_xp"`
	SolvedCount int    `json:"solved_count"`
}

type Leaderboard struct {
	Entries         []LeaderboardEntry `json:"entries"`
	Total           int                `json:"total"`
	CurrentUserRank int                `json:"current_user_rank"`
}
