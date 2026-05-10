CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    username      VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'student'
                  CHECK (role IN ('student', 'admin')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON auth.users(email);
CREATE INDEX idx_users_username ON auth.users(username);
