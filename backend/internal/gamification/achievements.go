package gamification

// CheckContext provides data needed by achievement checkers.
type CheckContext struct {
	TotalSolved              int
	HardTasksSolved          int
	TotalXP                  int
	ConsecutiveFirstTrySolves int
}

type AchievementChecker struct {
	Code  string
	Check func(ctx CheckContext) bool
}

var AchievementCheckers = []AchievementChecker{
	{
		Code: "first_solve",
		Check: func(ctx CheckContext) bool {
			return ctx.TotalSolved >= 1
		},
	},
	{
		Code: "streak_5",
		Check: func(ctx CheckContext) bool {
			return ctx.ConsecutiveFirstTrySolves >= 5
		},
	},
	{
		Code: "ten_tasks",
		Check: func(ctx CheckContext) bool {
			return ctx.TotalSolved >= 10
		},
	},
	{
		Code: "hard_solver",
		Check: func(ctx CheckContext) bool {
			return ctx.HardTasksSolved >= 5
		},
	},
	{
		Code: "century",
		Check: func(ctx CheckContext) bool {
			return ctx.TotalXP >= 100
		},
	},
}
