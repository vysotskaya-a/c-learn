package gamification

import (
	"context"
	"log"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/errs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AwardXP(ctx context.Context, req models.XPAwardRequest) (*models.XPAwardResponse, *errs.AppError) {
	// Idempotency check
	already, err := s.repo.XPAlreadyAwarded(ctx, req.SolutionID)
	if err != nil {
		return nil, errs.NewInternal("check idempotency failed")
	}
	if already {
		cached, err := s.repo.GetCachedXPAward(ctx, req.SolutionID)
		if err != nil {
			return nil, errs.NewInternal("get cached xp failed")
		}
		return cached, nil
	}

	// Task-level idempotency: award XP only on the first successful solve of a task.
	taskAlready, err := s.repo.XPAlreadyAwardedForTask(ctx, req.UserID, req.TaskID)
	if err != nil {
		return nil, errs.NewInternal("check task idempotency failed")
	}
	if taskAlready {
		totalXP, _ := s.repo.GetTotalXP(ctx, req.UserID)
		return &models.XPAwardResponse{
			XPAwarded:       0,
			TotalXP:         totalXP,
			Level:           CalculateLevel(totalXP),
			NewAchievements: []models.Achievement{},
		}, nil
	}

	// Calculate XP
	xpAwarded := CalculateXP(req.Difficulty, req.AttemptNumber)

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, errs.NewInternal("begin transaction failed")
	}
	defer tx.Rollback()

	// Insert XP log
	if err := s.repo.InsertXPLog(ctx, tx, req.UserID, req.SolutionID, req.TaskID, xpAwarded, req.Difficulty, req.AttemptNumber); err != nil {
		return nil, errs.NewInternal("insert xp log failed")
	}

	// Check achievement triggers
	totalXP, _ := s.repo.GetTotalXP(ctx, req.UserID)
	totalXP += xpAwarded // include current award

	solvedCount, _ := s.repo.GetSolvedCount(ctx, req.UserID)
	solvedCount++ // include current

	hardSolved, _ := s.repo.GetHardTasksSolved(ctx, req.UserID)
	if req.Difficulty == "hard" {
		hardSolved++
	}

	streak, _ := s.repo.GetConsecutiveFirstTrySolves(ctx, req.UserID)
	if req.AttemptNumber == 1 {
		streak++
	}

	checkCtx := CheckContext{
		TotalSolved:              solvedCount,
		HardTasksSolved:          hardSolved,
		TotalXP:                  totalXP,
		ConsecutiveFirstTrySolves: streak,
	}

	var newAchievements []models.Achievement
	for _, checker := range AchievementCheckers {
		if checker.Check(checkCtx) {
			achID, title, desc, iconURL, err := s.repo.GetAchievementByCode(ctx, checker.Code)
			if err != nil {
				log.Printf("achievement %s not found: %v", checker.Code, err)
				continue
			}

			hasIt, _ := s.repo.UserHasAchievement(ctx, req.UserID, achID)
			if hasIt {
				continue
			}

			if err := s.repo.AwardAchievement(ctx, tx, req.UserID, achID); err != nil {
				log.Printf("award achievement %s failed: %v", checker.Code, err)
				continue
			}

			newAchievements = append(newAchievements, models.Achievement{
				Code:    checker.Code,
				Title:   title,
				Description: desc,
				IconURL: iconURL,
			})
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errs.NewInternal("commit transaction failed")
	}

	if newAchievements == nil {
		newAchievements = []models.Achievement{}
	}

	return &models.XPAwardResponse{
		XPAwarded:       xpAwarded,
		TotalXP:         totalXP,
		Level:           CalculateLevel(totalXP),
		NewAchievements: newAchievements,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*models.Profile, *errs.AppError) {
	totalXP, err := s.repo.GetTotalXP(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal("get total xp failed")
	}

	solvedCount, err := s.repo.GetSolvedCount(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal("get solved count failed")
	}

	achievements, err := s.repo.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal("get achievements failed")
	}
	if achievements == nil {
		achievements = []models.UserAchievement{}
	}

	return &models.Profile{
		UserID:       userID,
		TotalXP:      totalXP,
		Level:        CalculateLevel(totalXP),
		SolvedCount:  solvedCount,
		Achievements: achievements,
	}, nil
}

func (s *Service) GetLeaderboard(ctx context.Context, userID string) (*models.Leaderboard, *errs.AppError) {
	entries, err := s.repo.GetLeaderboard(ctx, 50)
	if err != nil {
		return nil, errs.NewInternal("get leaderboard failed")
	}
	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	rank, _ := s.repo.GetUserRank(ctx, userID)

	total := len(entries)
	return &models.Leaderboard{
		Entries:         entries,
		Total:           total,
		CurrentUserRank: rank,
	}, nil
}
