CREATE SCHEMA IF NOT EXISTS lms;

-- Modules
CREATE TABLE lms.modules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_modules_sort ON lms.modules(sort_order);

-- Lessons
CREATE TABLE lms.lessons (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id   UUID NOT NULL REFERENCES lms.modules(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    theory_md   TEXT NOT NULL DEFAULT '',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_lessons_module ON lms.lessons(module_id);
CREATE INDEX idx_lessons_sort ON lms.lessons(module_id, sort_order);

-- Difficulty enum
CREATE TYPE lms.difficulty AS ENUM ('easy', 'medium', 'hard');

-- Tasks
CREATE TABLE lms.tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id   UUID NOT NULL REFERENCES lms.lessons(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty  lms.difficulty NOT NULL DEFAULT 'easy',
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_tasks_lesson ON lms.tasks(lesson_id);

-- Test Cases
CREATE TABLE lms.test_cases (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id     UUID NOT NULL REFERENCES lms.tasks(id) ON DELETE CASCADE,
    input       TEXT NOT NULL DEFAULT '',
    expected    TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    is_sample   BOOLEAN NOT NULL DEFAULT false
);
CREATE INDEX idx_test_cases_task ON lms.test_cases(task_id);

-- Verdict enum
CREATE TYPE lms.verdict AS ENUM ('ok', 'compilation_error', 'wrong_answer',
                                  'time_limit_exceeded', 'runtime_error');

-- Solutions
CREATE TABLE lms.solutions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    task_id         UUID NOT NULL REFERENCES lms.tasks(id) ON DELETE CASCADE,
    source_code     TEXT NOT NULL,
    verdict         lms.verdict NOT NULL,
    compiler_output TEXT,
    tests_passed    INTEGER NOT NULL DEFAULT 0,
    tests_total     INTEGER NOT NULL DEFAULT 0,
    exec_time_ms    INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_solutions_user ON lms.solutions(user_id);
CREATE INDEX idx_solutions_task ON lms.solutions(task_id);
CREATE INDEX idx_solutions_user_task ON lms.solutions(user_id, task_id);

-- Lesson Status enum
CREATE TYPE lms.lesson_status AS ENUM ('not_started', 'in_progress', 'completed');

-- User Lesson Progress
CREATE TABLE lms.user_lesson_progress (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    lesson_id   UUID NOT NULL REFERENCES lms.lessons(id) ON DELETE CASCADE,
    status      lms.lesson_status NOT NULL DEFAULT 'not_started',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, lesson_id)
);
CREATE INDEX idx_progress_user ON lms.user_lesson_progress(user_id);
