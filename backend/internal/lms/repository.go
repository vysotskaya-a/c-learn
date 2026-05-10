package lms

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/c-learn/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ---- Modules ----

func (r *Repository) ListModules(ctx context.Context) ([]models.Module, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, COALESCE(description,''), sort_order, created_at, updated_at
		 FROM lms.modules ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *Repository) CreateModule(ctx context.Context, title, description string, sortOrder int) (*models.Module, error) {
	m := &models.Module{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO lms.modules (title, description, sort_order) VALUES ($1, $2, $3)
		 RETURNING id, title, description, sort_order, created_at, updated_at`,
		title, description, sortOrder,
	).Scan(&m.ID, &m.Title, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *Repository) UpdateModule(ctx context.Context, id, title, description string, sortOrder int) (*models.Module, error) {
	m := &models.Module{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE lms.modules SET title=$2, description=$3, sort_order=$4, updated_at=now()
		 WHERE id=$1 RETURNING id, title, description, sort_order, created_at, updated_at`,
		id, title, description, sortOrder,
	).Scan(&m.ID, &m.Title, &m.Description, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *Repository) DeleteModule(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM lms.modules WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ---- Lessons ----

func (r *Repository) ListLessonsByModule(ctx context.Context, moduleID string) ([]models.Lesson, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, module_id, title, theory_md, sort_order, created_at, updated_at
		 FROM lms.lessons WHERE module_id=$1 ORDER BY sort_order`, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		if err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.TheoryMD, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *Repository) GetLesson(ctx context.Context, id string) (*models.Lesson, error) {
	l := &models.Lesson{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, module_id, title, theory_md, sort_order, created_at, updated_at
		 FROM lms.lessons WHERE id=$1`, id,
	).Scan(&l.ID, &l.ModuleID, &l.Title, &l.TheoryMD, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return l, err
}

func (r *Repository) CreateLesson(ctx context.Context, moduleID, title, theoryMD string, sortOrder int) (*models.Lesson, error) {
	l := &models.Lesson{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO lms.lessons (module_id, title, theory_md, sort_order)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, module_id, title, theory_md, sort_order, created_at, updated_at`,
		moduleID, title, theoryMD, sortOrder,
	).Scan(&l.ID, &l.ModuleID, &l.Title, &l.TheoryMD, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	return l, err
}

func (r *Repository) UpdateLesson(ctx context.Context, id, title, theoryMD string, sortOrder int) (*models.Lesson, error) {
	l := &models.Lesson{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE lms.lessons SET title=$2, theory_md=$3, sort_order=$4, updated_at=now()
		 WHERE id=$1 RETURNING id, module_id, title, theory_md, sort_order, created_at, updated_at`,
		id, title, theoryMD, sortOrder,
	).Scan(&l.ID, &l.ModuleID, &l.Title, &l.TheoryMD, &l.SortOrder, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return l, err
}

func (r *Repository) DeleteLesson(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM lms.lessons WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ---- Tasks ----

func (r *Repository) ListTasksByLesson(ctx context.Context, lessonID string) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, lesson_id, title, description, difficulty, sort_order, created_at, updated_at
		 FROM lms.tasks WHERE lesson_id=$1 ORDER BY sort_order`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.Difficulty, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) GetTask(ctx context.Context, id string) (*models.Task, error) {
	t := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, lesson_id, title, description, difficulty, sort_order, created_at, updated_at
		 FROM lms.tasks WHERE id=$1`, id,
	).Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.Difficulty, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *Repository) CreateTask(ctx context.Context, lessonID, title, description, difficulty string, sortOrder int) (*models.Task, error) {
	t := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO lms.tasks (lesson_id, title, description, difficulty, sort_order)
		 VALUES ($1, $2, $3, $4::lms.difficulty, $5)
		 RETURNING id, lesson_id, title, description, difficulty, sort_order, created_at, updated_at`,
		lessonID, title, description, difficulty, sortOrder,
	).Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.Difficulty, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *Repository) UpdateTask(ctx context.Context, id, title, description, difficulty string, sortOrder int) (*models.Task, error) {
	t := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE lms.tasks SET title=$2, description=$3, difficulty=$4::lms.difficulty, sort_order=$5, updated_at=now()
		 WHERE id=$1 RETURNING id, lesson_id, title, description, difficulty, sort_order, created_at, updated_at`,
		id, title, description, difficulty, sortOrder,
	).Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.Difficulty, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *Repository) DeleteTask(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM lms.tasks WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ---- Test Cases ----

func (r *Repository) GetTestCasesByTask(ctx context.Context, taskID string) ([]models.TestCase, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, task_id, input, expected, sort_order, is_sample
		 FROM lms.test_cases WHERE task_id=$1 ORDER BY sort_order`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tcs []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.TaskID, &tc.Input, &tc.Expected, &tc.SortOrder, &tc.IsSample); err != nil {
			return nil, err
		}
		tcs = append(tcs, tc)
	}
	return tcs, nil
}

func (r *Repository) GetSampleTestCases(ctx context.Context, taskID string) ([]models.SampleTest, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT input, expected FROM lms.test_cases
		 WHERE task_id=$1 AND is_sample=true ORDER BY sort_order`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var samples []models.SampleTest
	for rows.Next() {
		var s models.SampleTest
		if err := rows.Scan(&s.Input, &s.Expected); err != nil {
			return nil, err
		}
		samples = append(samples, s)
	}
	return samples, nil
}

func (r *Repository) CreateTestCase(ctx context.Context, taskID, input, expected string, sortOrder int, isSample bool) (*models.TestCase, error) {
	tc := &models.TestCase{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO lms.test_cases (task_id, input, expected, sort_order, is_sample)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, task_id, input, expected, sort_order, is_sample`,
		taskID, input, expected, sortOrder, isSample,
	).Scan(&tc.ID, &tc.TaskID, &tc.Input, &tc.Expected, &tc.SortOrder, &tc.IsSample)
	return tc, err
}

func (r *Repository) DeleteTestCasesByTask(ctx context.Context, taskID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lms.test_cases WHERE task_id=$1`, taskID)
	return err
}

// ---- Solutions ----

func (r *Repository) CreateSolutionTx(ctx context.Context, tx *sql.Tx, userID, taskID, sourceCode, verdict, compilerOutput string, testsPassed, testsTotal, execTimeMs int) (string, error) {
	var id string
	err := tx.QueryRowContext(ctx,
		`INSERT INTO lms.solutions (user_id, task_id, source_code, verdict, compiler_output, tests_passed, tests_total, exec_time_ms)
		 VALUES ($1, $2, $3, $4::lms.verdict, $5, $6, $7, $8) RETURNING id`,
		userID, taskID, sourceCode, verdict, compilerOutput, testsPassed, testsTotal, execTimeMs,
	).Scan(&id)
	return id, err
}

func (r *Repository) CountAttempts(ctx context.Context, userID, taskID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM lms.solutions WHERE user_id=$1 AND task_id=$2`,
		userID, taskID,
	).Scan(&count)
	return count, err
}

func (r *Repository) HasSolvedTask(ctx context.Context, userID, taskID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM lms.solutions WHERE user_id=$1 AND task_id=$2 AND verdict='ok')`,
		userID, taskID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) ListSolutions(ctx context.Context, userID string, limit, offset int) ([]models.Solution, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, task_id, source_code, verdict, COALESCE(compiler_output,''),
		        tests_passed, tests_total, COALESCE(exec_time_ms,0), created_at
		 FROM lms.solutions WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sols []models.Solution
	for rows.Next() {
		var s models.Solution
		if err := rows.Scan(&s.ID, &s.UserID, &s.TaskID, &s.SourceCode, &s.Verdict, &s.CompilerOutput,
			&s.TestsPassed, &s.TestsTotal, &s.ExecTimeMs, &s.CreatedAt); err != nil {
			return nil, err
		}
		sols = append(sols, s)
	}
	return sols, nil
}

// ---- Progress ----

func (r *Repository) GetLessonProgress(ctx context.Context, userID, lessonID string) (string, error) {
	var status string
	err := r.db.QueryRowContext(ctx,
		`SELECT status FROM lms.user_lesson_progress WHERE user_id=$1 AND lesson_id=$2`,
		userID, lessonID,
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "not_started", nil
	}
	return status, err
}

func (r *Repository) GetAllProgress(ctx context.Context, userID string) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT lesson_id, status FROM lms.user_lesson_progress WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]string)
	for rows.Next() {
		var lid, status string
		if err := rows.Scan(&lid, &status); err != nil {
			return nil, err
		}
		m[lid] = status
	}
	return m, nil
}

func (r *Repository) UpsertProgress(ctx context.Context, userID, lessonID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO lms.user_lesson_progress (user_id, lesson_id, status, updated_at)
		 VALUES ($1, $2, $3::lms.lesson_status, now())
		 ON CONFLICT (user_id, lesson_id) DO UPDATE SET status=$3::lms.lesson_status, updated_at=now()`,
		userID, lessonID, status,
	)
	return err
}

func (r *Repository) UpsertProgressTx(ctx context.Context, tx *sql.Tx, userID, lessonID, status string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO lms.user_lesson_progress (user_id, lesson_id, status, updated_at)
		 VALUES ($1, $2, $3::lms.lesson_status, now())
		 ON CONFLICT (user_id, lesson_id) DO UPDATE SET status=$3::lms.lesson_status, updated_at=now()`,
		userID, lessonID, status,
	)
	return err
}

// Check if all tasks in a lesson are solved by user
func (r *Repository) AllTasksSolved(ctx context.Context, tx *sql.Tx, userID, lessonID string) (bool, error) {
	var unsolved int
	err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM lms.tasks t
		 WHERE t.lesson_id = $1
		 AND NOT EXISTS (
		   SELECT 1 FROM lms.solutions s
		   WHERE s.task_id = t.id AND s.user_id = $2 AND s.verdict = 'ok'
		 )`, lessonID, userID,
	).Scan(&unsolved)
	if err != nil {
		return false, err
	}
	return unsolved == 0, nil
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// Get lesson_id for a task
func (r *Repository) GetTaskLessonID(ctx context.Context, taskID string) (string, error) {
	var lessonID string
	err := r.db.QueryRowContext(ctx,
		`SELECT lesson_id FROM lms.tasks WHERE id=$1`, taskID,
	).Scan(&lessonID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return lessonID, err
}

// GetSolvedTaskIDsForLesson returns set of task IDs solved by user in a given lesson.
func (r *Repository) GetSolvedTaskIDsForLesson(ctx context.Context, userID, lessonID string) (map[string]bool, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT s.task_id FROM lms.solutions s
		 JOIN lms.tasks t ON t.id = s.task_id
		 WHERE s.user_id=$1 AND t.lesson_id=$2 AND s.verdict='ok'`,
		userID, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	solved := make(map[string]bool)
	for rows.Next() {
		var tid string
		if err := rows.Scan(&tid); err != nil {
			return nil, err
		}
		solved[tid] = true
	}
	return solved, nil
}
