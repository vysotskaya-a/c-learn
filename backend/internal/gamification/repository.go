package gamification

import (
	"context"
	"database/sql"

	"github.com/c-learn/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ---- XP Log ----

func (r *Repository) XPAlreadyAwardedForTask(ctx context.Context, userID, taskID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM gamification.user_xp_log WHERE user_id=$1 AND task_id=$2)`,
		userID, taskID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) XPAlreadyAwarded(ctx context.Context, solutionID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM gamification.user_xp_log WHERE solution_id=$1)`,
		solutionID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) GetCachedXPAward(ctx context.Context, solutionID string) (*models.XPAwardResponse, error) {
	var xpAwarded int
	var userID, difficulty string
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, xp_awarded, difficulty FROM gamification.user_xp_log WHERE solution_id=$1`,
		solutionID,
	).Scan(&userID, &xpAwarded, &difficulty)
	if err != nil {
		return nil, err
	}

	totalXP, _ := r.GetTotalXP(ctx, userID)
	return &models.XPAwardResponse{
		XPAwarded:        xpAwarded,
		TotalXP:          totalXP,
		Level:            CalculateLevel(totalXP),
		NewAchievements:  []models.Achievement{},
		AlreadyProcessed: true,
	}, nil
}

func (r *Repository) InsertXPLog(ctx context.Context, tx *sql.Tx, userID, solutionID, taskID string, xpAwarded int, difficulty string, attemptNumber int) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO gamification.user_xp_log (user_id, solution_id, task_id, xp_awarded, difficulty, attempt_number)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, solutionID, taskID, xpAwarded, difficulty, attemptNumber,
	)
	return err
}

func (r *Repository) GetTotalXP(ctx context.Context, userID string) (int, error) {
	var total sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT SUM(xp_awarded) FROM gamification.user_xp_log WHERE user_id=$1`, userID,
	).Scan(&total)
	if err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return int(total.Int64), nil
}

func (r *Repository) GetSolvedCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT task_id) FROM gamification.user_xp_log WHERE user_id=$1`, userID,
	).Scan(&count)
	return count, err
}

func (r *Repository) GetHardTasksSolved(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT task_id) FROM gamification.user_xp_log
		 WHERE user_id=$1 AND difficulty='hard'`, userID,
	).Scan(&count)
	return count, err
}

func (r *Repository) GetConsecutiveFirstTrySolves(ctx context.Context, userID string) (int, error) {
	// Count consecutive solutions where attempt_number = 1, from most recent
	rows, err := r.db.QueryContext(ctx,
		`SELECT attempt_number FROM gamification.user_xp_log
		 WHERE user_id=$1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	streak := 0
	for rows.Next() {
		var attempt int
		if err := rows.Scan(&attempt); err != nil {
			return streak, err
		}
		if attempt == 1 {
			streak++
		} else {
			break
		}
	}
	return streak, nil
}

// ---- Achievements ----

func (r *Repository) GetAchievementByCode(ctx context.Context, code string) (id, title, description, iconURL string, err error) {
	var iconNull sql.NullString
	err = r.db.QueryRowContext(ctx,
		`SELECT id, title, description, COALESCE(icon_url, '') FROM gamification.achievements WHERE code=$1`, code,
	).Scan(&id, &title, &description, &iconNull)
	if iconNull.Valid {
		iconURL = iconNull.String
	}
	return
}

func (r *Repository) UserHasAchievement(ctx context.Context, userID, achievementID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM gamification.user_achievements WHERE user_id=$1 AND achievement_id=$2)`,
		userID, achievementID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) AwardAchievement(ctx context.Context, tx *sql.Tx, userID, achievementID string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO gamification.user_achievements (user_id, achievement_id)
		 VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, achievementID,
	)
	return err
}

func (r *Repository) GetUserAchievements(ctx context.Context, userID string) ([]models.UserAchievement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT a.code, a.title, a.description, COALESCE(a.icon_url, ''), ua.awarded_at
		 FROM gamification.user_achievements ua
		 JOIN gamification.achievements a ON a.id = ua.achievement_id
		 WHERE ua.user_id=$1
		 ORDER BY ua.awarded_at`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []models.UserAchievement
	for rows.Next() {
		var a models.UserAchievement
		if err := rows.Scan(&a.Code, &a.Title, &a.Description, &a.IconURL, &a.AwardedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

// ---- Leaderboard ----

func (r *Repository) GetLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT l.user_id, COALESCE(u.username, '') AS username,
		        SUM(l.xp_awarded) AS total_xp,
		        COUNT(DISTINCT l.task_id) AS solved_count
		 FROM gamification.user_xp_log l
		 LEFT JOIN auth.users u ON u.id = l.user_id
		 GROUP BY l.user_id, u.username
		 ORDER BY total_xp DESC
		 LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e models.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.TotalXP, &e.SolvedCount); err != nil {
			return nil, err
		}
		e.Rank = rank
		rank++
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *Repository) GetUserRank(ctx context.Context, userID string) (int, error) {
	var rank int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) + 1 FROM (
			SELECT user_id, SUM(xp_awarded) as total_xp
			FROM gamification.user_xp_log
			GROUP BY user_id
		) sub WHERE sub.total_xp > COALESCE(
			(SELECT SUM(xp_awarded) FROM gamification.user_xp_log WHERE user_id=$1), 0
		)`, userID,
	).Scan(&rank)
	return rank, err
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
