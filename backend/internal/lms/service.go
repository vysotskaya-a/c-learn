package lms

import (
	"context"
	"log"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/errs"
	"github.com/c-learn/pkg/validator"
)

type Service struct {
	repo   *Repository
	runner *RunnerClient
	gamif  *GamificationClient
}

func NewService(repo *Repository, runner *RunnerClient, gamif *GamificationClient) *Service {
	return &Service{repo: repo, runner: runner, gamif: gamif}
}

// ---- Course Tree ----

func (s *Service) GetCourseTree(ctx context.Context, userID string) ([]models.ModuleWithLessons, *errs.AppError) {
	modules, err := s.repo.ListModules(ctx)
	if err != nil {
		return nil, errs.NewInternal("list modules failed")
	}

	progress, err := s.repo.GetAllProgress(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal("get progress failed")
	}

	var result []models.ModuleWithLessons
	for _, m := range modules {
		lessons, err := s.repo.ListLessonsByModule(ctx, m.ID)
		if err != nil {
			return nil, errs.NewInternal("list lessons failed")
		}
		var briefs []models.LessonBrief
		for _, l := range lessons {
			status := "not_started"
			if s, ok := progress[l.ID]; ok {
				status = s
			}
			briefs = append(briefs, models.LessonBrief{
				ID:        l.ID,
				Title:     l.Title,
				SortOrder: l.SortOrder,
				Status:    status,
			})
		}
		if briefs == nil {
			briefs = []models.LessonBrief{}
		}
		result = append(result, models.ModuleWithLessons{
			ID:        m.ID,
			Title:     m.Title,
			SortOrder: m.SortOrder,
			Lessons:   briefs,
		})
	}
	if result == nil {
		result = []models.ModuleWithLessons{}
	}
	return result, nil
}

// ---- Lesson Detail ----

func (s *Service) GetLesson(ctx context.Context, lessonID, userID string) (*models.LessonDetail, *errs.AppError) {
	lesson, err := s.repo.GetLesson(ctx, lessonID)
	if err != nil {
		return nil, errs.NewInternal("get lesson failed")
	}
	if lesson == nil {
		return nil, errs.NewNotFound("lesson not found")
	}

	tasks, err := s.repo.ListTasksByLesson(ctx, lessonID)
	if err != nil {
		return nil, errs.NewInternal("list tasks failed")
	}

	solvedMap, err := s.repo.GetSolvedTaskIDsForLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, errs.NewInternal("get solved tasks failed")
	}

	// Mark progress as in_progress if not_started
	status, _ := s.repo.GetLessonProgress(ctx, userID, lessonID)
	if status == "not_started" {
		_ = s.repo.UpsertProgress(ctx, userID, lessonID, "in_progress")
		status = "in_progress"
	}

	var taskDetails []models.TaskDetail
	for _, t := range tasks {
		samples, _ := s.repo.GetSampleTestCases(ctx, t.ID)
		if samples == nil {
			samples = []models.SampleTest{}
		}
		taskDetails = append(taskDetails, models.TaskDetail{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Difficulty:  t.Difficulty,
			IsSolved:    solvedMap[t.ID],
			Samples:     samples,
		})
	}
	if taskDetails == nil {
		taskDetails = []models.TaskDetail{}
	}

	return &models.LessonDetail{
		ID:       lesson.ID,
		ModuleID: lesson.ModuleID,
		Title:    lesson.Title,
		TheoryMD: lesson.TheoryMD,
		Status:   status,
		Tasks:    taskDetails,
	}, nil
}

// ---- Run (test run, no save) ----

func (s *Service) RunCode(ctx context.Context, taskID, sourceCode, stdin string) (*models.RunResponse, *errs.AppError) {
	if err := validator.ValidateSourceCode(sourceCode); err != nil {
		return nil, errs.NewValidation(err.Error(), map[string]any{"field": "source_code"})
	}

	// Verify task exists
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil || task == nil {
		return nil, errs.NewNotFound("task not found")
	}

	result, err := s.runner.Run(ctx, models.RunRequest{
		SourceCode: sourceCode,
		TestCases: []models.RunTestCase{
			{Input: stdin, Expected: ""},
		},
		Mode: "run",
	})
	if err != nil {
		log.Printf("runner error: %v", err)
		return nil, errs.NewInternal("code runner unavailable")
	}

	return &models.RunResponse{
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
		ExitCode:   0,
		ExecTimeMs: result.ExecTimeMs,
	}, nil
}

// ---- Submit (full flow with orchestration) ----

func (s *Service) Submit(ctx context.Context, userID, taskID, sourceCode string) (*models.SubmitResponse, *errs.AppError) {
	if err := validator.ValidateSourceCode(sourceCode); err != nil {
		return nil, errs.NewValidation(err.Error(), map[string]any{"field": "source_code"})
	}

	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil || task == nil {
		return nil, errs.NewNotFound("task not found")
	}

	// Get test cases
	testCases, err := s.repo.GetTestCasesByTask(ctx, taskID)
	if err != nil {
		return nil, errs.NewInternal("get test cases failed")
	}

	var runTCs []models.RunTestCase
	for _, tc := range testCases {
		runTCs = append(runTCs, models.RunTestCase{Input: tc.Input, Expected: tc.Expected})
	}

	// Call Code Runner
	result, err := s.runner.Run(ctx, models.RunRequest{
		SourceCode: sourceCode,
		TestCases:  runTCs,
		Mode:       "judge",
	})
	if err != nil {
		log.Printf("runner error: %v", err)
		return nil, errs.NewInternal("code runner unavailable")
	}

	// Count attempts
	attempts, _ := s.repo.CountAttempts(ctx, userID, taskID)
	attemptNumber := attempts + 1

	// Begin transaction in lms_db
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, errs.NewInternal("begin transaction failed")
	}
	defer tx.Rollback()

	solutionID, err := s.repo.CreateSolutionTx(ctx, tx, userID, taskID, sourceCode,
		result.Verdict, result.CompilerOutput, result.TestsPassed, result.TestsTotal, result.ExecTimeMs)
	if err != nil {
		return nil, errs.NewInternal("save solution failed")
	}

	// If OK, check if lesson is complete
	lessonID := task.LessonID
	if result.Verdict == "ok" {
		allSolved, err := s.repo.AllTasksSolved(ctx, tx, userID, lessonID)
		if err == nil && allSolved {
			_ = s.repo.UpsertProgressTx(ctx, tx, userID, lessonID, "completed")
		} else if result.Verdict == "ok" {
			_ = s.repo.UpsertProgressTx(ctx, tx, userID, lessonID, "in_progress")
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errs.NewInternal("commit transaction failed")
	}

	// Award XP if verdict is OK
	resp := &models.SubmitResponse{
		SolutionID:           solutionID,
		Verdict:              result.Verdict,
		TestsPassed:          result.TestsPassed,
		TestsTotal:           result.TestsTotal,
		ExecTimeMs:           result.ExecTimeMs,
		CompilerOutput:       result.CompilerOutput,
		XPAwarded:            0,
		AchievementsUnlocked: []models.Achievement{},
	}

	if result.FailedTest != nil {
		resp.FailedTest = *result.FailedTest
	}

	if result.Verdict == "ok" {
		xpResp, err := s.gamif.AwardXP(ctx, models.XPAwardRequest{
			UserID:        userID,
			SolutionID:    solutionID,
			TaskID:        taskID,
			Difficulty:    task.Difficulty,
			AttemptNumber: attemptNumber,
		})
		if err != nil {
			log.Printf("gamification error (non-fatal): %v", err)
			// Graceful degradation: solution saved, XP not awarded
		} else {
			resp.XPAwarded = xpResp.XPAwarded
			if xpResp.NewAchievements != nil {
				resp.AchievementsUnlocked = xpResp.NewAchievements
			}
		}
	}

	return resp, nil
}

// ---- Solutions ----

func (s *Service) ListSolutions(ctx context.Context, userID string, limit, offset int) ([]models.Solution, *errs.AppError) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	sols, err := s.repo.ListSolutions(ctx, userID, limit, offset)
	if err != nil {
		return nil, errs.NewInternal("list solutions failed")
	}
	if sols == nil {
		sols = []models.Solution{}
	}
	return sols, nil
}
