CREATE SCHEMA IF NOT EXISTS gamification;

-- Achievements catalog
CREATE TABLE gamification.achievements (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        VARCHAR(100) NOT NULL UNIQUE,
    title       VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon_url    VARCHAR(500),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User achievements
CREATE TABLE gamification.user_achievements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    achievement_id  UUID NOT NULL REFERENCES gamification.achievements(id) ON DELETE CASCADE,
    awarded_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, achievement_id)
);
CREATE INDEX idx_user_achievements_user ON gamification.user_achievements(user_id);

-- XP log
CREATE TABLE gamification.user_xp_log (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL,
    solution_id    UUID NOT NULL UNIQUE,
    task_id        UUID NOT NULL,
    xp_awarded     INTEGER NOT NULL,
    difficulty     VARCHAR(20) NOT NULL,
    attempt_number INTEGER NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_xp_log_user ON gamification.user_xp_log(user_id);
CREATE INDEX idx_xp_log_solution ON gamification.user_xp_log(solution_id);

-- Seed achievements
INSERT INTO gamification.achievements (code, title, description, icon_url) VALUES
('first_solve',  'Первая решённая задача', 'Решите свою первую задачу', '/badges/first_solve.svg'),
('streak_5',     'Серия из 5',             '5 задач подряд с первой попытки', '/badges/streak_5.svg'),
('ten_tasks',    '10 задач решено',        'Решите 10 задач', '/badges/ten_tasks.svg'),
('hard_solver',  'Покоритель сложных',     'Решите 5 сложных задач', '/badges/hard_solver.svg'),
('century',      'Сотня XP',              'Наберите 100 XP', '/badges/century.svg');
