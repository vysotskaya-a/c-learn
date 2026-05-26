package gamification

import "testing"

func TestAchievementCheckers(t *testing.T) {
	tests := []struct {
		code    string
		ctx     CheckContext
		unlock  bool
	}{
		{code: "first_solve", ctx: CheckContext{TotalSolved: 1}, unlock: true},
		{code: "first_solve", ctx: CheckContext{TotalSolved: 0}, unlock: false},
		{code: "streak_5", ctx: CheckContext{ConsecutiveFirstTrySolves: 5}, unlock: true},
		{code: "streak_5", ctx: CheckContext{ConsecutiveFirstTrySolves: 4}, unlock: false},
		{code: "ten_tasks", ctx: CheckContext{TotalSolved: 10}, unlock: true},
		{code: "hard_solver", ctx: CheckContext{HardTasksSolved: 5}, unlock: true},
		{code: "century", ctx: CheckContext{TotalXP: 100}, unlock: true},
		{code: "century", ctx: CheckContext{TotalXP: 99}, unlock: false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			var checker AchievementChecker
			for _, c := range AchievementCheckers {
				if c.Code == tt.code {
					checker = c
					break
				}
			}
			if checker.Code == "" {
				t.Fatalf("checker %s not found", tt.code)
			}
			got := checker.Check(tt.ctx)
			if got != tt.unlock {
				t.Fatalf("Check(%+v) = %v, want %v", tt.ctx, got, tt.unlock)
			}
		})
	}
}
