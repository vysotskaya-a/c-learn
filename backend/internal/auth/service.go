package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/errs"
	"github.com/c-learn/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, email, username, passwordHash, role string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (id, userID string, expiresAt time.Time, err error)
	DeleteRefreshToken(ctx context.Context, id string) error
}

type Service struct {
	repo userRepository
	jwt  *JWTManager
}

func NewService(repo *Repository, jwt *JWTManager) *Service {
	return &Service{repo: repo, jwt: jwt}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*models.User, *errs.AppError) {
	if err := validator.ValidateEmail(req.Email); err != nil {
		return nil, errs.NewValidation(err.Error(), map[string]any{"field": "email"})
	}
	if err := validator.ValidateUsername(req.Username); err != nil {
		return nil, errs.NewValidation(err.Error(), map[string]any{"field": "username"})
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, errs.NewValidation(err.Error(), map[string]any{"field": "password", "rule": "min_length", "min": 8})
	}

	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errs.NewInternal("database error")
	}
	if exists {
		return nil, errs.NewConflict("email already taken")
	}

	exists, err = s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, errs.NewInternal("database error")
	}
	if exists {
		return nil, errs.NewConflict("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errs.NewInternal("hash password failed")
	}

	user, err := s.repo.CreateUser(ctx, req.Email, req.Username, string(hash), "student")
	if err != nil {
		return nil, errs.NewInternal("create user failed")
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, *errs.AppError) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.NewInternal("database error")
	}
	if user == nil {
		return nil, errs.NewUnauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errs.NewUnauthorized("invalid credentials")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, errs.NewInternal("generate access token failed")
	}

	rawRefresh, hashRefresh, expiresAt, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, errs.NewInternal("generate refresh token failed")
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, hashRefresh, expiresAt); err != nil {
		return nil, errs.NewInternal("save refresh token failed")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*TokenResponse, *errs.AppError) {
	hash := HashRefreshToken(req.RefreshToken)
	id, userID, expiresAt, err := s.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return nil, errs.NewInternal("database error")
	}
	if id == "" {
		return nil, errs.NewUnauthorized("invalid refresh token")
	}
	if time.Now().After(expiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, id)
		return nil, errs.NewUnauthorized("refresh token expired")
	}

	// Rotate: delete old, create new
	_ = s.repo.DeleteRefreshToken(ctx, id)

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errs.NewUnauthorized("user not found")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, errs.NewInternal("generate access token failed")
	}

	rawRefresh, hashRefresh, newExpires, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, errs.NewInternal("generate refresh token failed")
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, hashRefresh, newExpires); err != nil {
		return nil, errs.NewInternal("save refresh token failed")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}

func (s *Service) GetMe(ctx context.Context, userID string) (*models.User, *errs.AppError) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal("database error")
	}
	if user == nil {
		return nil, errs.NewNotFound("user not found")
	}
	return user, nil
}

func (s *Service) GetUserInfo(ctx context.Context, userID string) (*models.UserInfo, *errs.AppError) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errs.NewInternal(fmt.Sprintf("database error: %v", err))
	}
	if user == nil {
		return nil, errs.NewNotFound("user not found")
	}
	return &models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
