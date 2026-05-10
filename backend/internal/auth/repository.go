package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/c-learn/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, email, username, passwordHash, role string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO auth.users (email, username, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, email, username, role, created_at, updated_at`,
		email, username, passwordHash, role,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, role, created_at, updated_at
		 FROM auth.users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, role, created_at, updated_at
		 FROM auth.users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM auth.users WHERE email = $1)`, email,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM auth.users WHERE username = $1)`, username,
	).Scan(&exists)
	return exists, err
}

// Refresh tokens

func (r *Repository) SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth.refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)`, userID, tokenHash, expiresAt,
	)
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (id, userID string, expiresAt time.Time, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT id, user_id, expires_at FROM auth.refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&id, &userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", "", time.Time{}, nil
	}
	return
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth.refresh_tokens WHERE id = $1`, id)
	return err
}

func (r *Repository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth.refresh_tokens WHERE user_id = $1`, userID)
	return err
}
