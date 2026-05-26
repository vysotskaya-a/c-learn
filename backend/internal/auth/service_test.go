package auth

import (
	"context"
	"testing"
	"time"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/errs"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	emails      map[string]bool
	usernames   map[string]bool
	usersByMail map[string]*models.User
	usersByID   map[string]*models.User
	refresh     map[string]refreshRecord
}

type refreshRecord struct {
	id        string
	userID    string
	expiresAt time.Time
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		emails:      map[string]bool{},
		usernames:   map[string]bool{},
		usersByMail: map[string]*models.User{},
		usersByID:   map[string]*models.User{},
		refresh:     map[string]refreshRecord{},
	}
}

func (m *mockUserRepo) EmailExists(_ context.Context, email string) (bool, error) {
	return m.emails[email], nil
}

func (m *mockUserRepo) UsernameExists(_ context.Context, username string) (bool, error) {
	return m.usernames[username], nil
}

func (m *mockUserRepo) CreateUser(_ context.Context, email, username, passwordHash, role string) (*models.User, error) {
	user := &models.User{
		ID:           "user-" + username,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
	}
	m.emails[email] = true
	m.usernames[username] = true
	m.usersByMail[email] = user
	m.usersByID[user.ID] = user
	return user, nil
}

func (m *mockUserRepo) GetUserByEmail(_ context.Context, email string) (*models.User, error) {
	u, ok := m.usersByMail[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) GetUserByID(_ context.Context, id string) (*models.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) SaveRefreshToken(_ context.Context, userID, tokenHash string, expiresAt time.Time) error {
	m.refresh[tokenHash] = refreshRecord{id: "rt-1", userID: userID, expiresAt: expiresAt}
	return nil
}

func (m *mockUserRepo) GetRefreshToken(_ context.Context, tokenHash string) (string, string, time.Time, error) {
	rec, ok := m.refresh[tokenHash]
	if !ok {
		return "", "", time.Time{}, nil
	}
	return rec.id, rec.userID, rec.expiresAt, nil
}

func (m *mockUserRepo) DeleteRefreshToken(_ context.Context, id string) error {
	for hash, rec := range m.refresh {
		if rec.id == id {
			delete(m.refresh, hash)
			break
		}
	}
	return nil
}

func testService(repo *mockUserRepo) *Service {
	return &Service{repo: repo, jwt: NewJWTManager("test-secret")}
}

func TestService_Register_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := testService(repo)

	user, appErr := svc.Register(context.Background(), RegisterRequest{
		Email:    "new@example.com",
		Username: "newuser",
		Password: "password1",
	})
	if appErr != nil {
		t.Fatalf("Register: %v", appErr)
	}
	if user.Email != "new@example.com" || user.Role != "student" {
		t.Fatalf("user = %+v", user)
	}
}

func TestService_Register_ValidationErrors(t *testing.T) {
	repo := newMockUserRepo()
	svc := testService(repo)

	_, appErr := svc.Register(context.Background(), RegisterRequest{
		Email: "bad", Username: "ab", Password: "short",
	})
	if appErr == nil || appErr.Code != errs.NewValidation("", nil).Code {
		t.Fatalf("expected validation error, got %v", appErr)
	}
}

func TestService_Register_ConflictEmail(t *testing.T) {
	repo := newMockUserRepo()
	repo.emails["taken@example.com"] = true
	svc := testService(repo)

	_, appErr := svc.Register(context.Background(), RegisterRequest{
		Email: "taken@example.com", Username: "user1", Password: "password1",
	})
	if appErr == nil || appErr.Code != errs.NewConflict("").Code {
		t.Fatalf("expected conflict, got %v", appErr)
	}
}

func TestService_Register_ConflictUsername(t *testing.T) {
	repo := newMockUserRepo()
	repo.usernames["taken"] = true
	svc := testService(repo)

	_, appErr := svc.Register(context.Background(), RegisterRequest{
		Email: "free@example.com", Username: "taken", Password: "password1",
	})
	if appErr == nil || appErr.Code != errs.NewConflict("").Code {
		t.Fatalf("expected conflict, got %v", appErr)
	}
}

func TestService_Login_Success(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), 10)
	repo.usersByMail["user@example.com"] = &models.User{
		ID: "u1", Email: "user@example.com", Username: "user", PasswordHash: string(hash), Role: "student",
	}
	svc := testService(repo)

	tokens, appErr := svc.Login(context.Background(), LoginRequest{
		Email: "user@example.com", Password: "password1",
	})
	if appErr != nil {
		t.Fatalf("Login: %v", appErr)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.TokenType != "Bearer" || tokens.ExpiresIn != 900 {
		t.Fatalf("tokens = %+v", tokens)
	}
}

func TestService_Login_InvalidCredentials(t *testing.T) {
	repo := newMockUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), 10)
	repo.usersByMail["user@example.com"] = &models.User{
		ID: "u1", Email: "user@example.com", PasswordHash: string(hash), Role: "student",
	}
	svc := testService(repo)

	_, appErr := svc.Login(context.Background(), LoginRequest{
		Email: "user@example.com", Password: "wrong",
	})
	if appErr == nil || appErr.Code != errs.NewUnauthorized("").Code {
		t.Fatalf("expected unauthorized, got %v", appErr)
	}

	_, appErr = svc.Login(context.Background(), LoginRequest{
		Email: "missing@example.com", Password: "password1",
	})
	if appErr == nil || appErr.Code != errs.NewUnauthorized("").Code {
		t.Fatalf("expected unauthorized for missing user, got %v", appErr)
	}
}

func TestService_Refresh_Success(t *testing.T) {
	repo := newMockUserRepo()
	repo.usersByID["u1"] = &models.User{ID: "u1", Role: "student"}
	svc := testService(repo)

	raw, hash, expiresAt, _ := svc.jwt.GenerateRefreshToken()
	repo.refresh[hash] = refreshRecord{id: "rt-1", userID: "u1", expiresAt: expiresAt}

	tokens, appErr := svc.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
	if appErr != nil {
		t.Fatalf("Refresh: %v", appErr)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("expected new token pair")
	}
}

func TestService_Refresh_Expired(t *testing.T) {
	repo := newMockUserRepo()
	repo.usersByID["u1"] = &models.User{ID: "u1", Role: "student"}
	svc := testService(repo)

	raw, hash, _, _ := svc.jwt.GenerateRefreshToken()
	repo.refresh[hash] = refreshRecord{id: "rt-1", userID: "u1", expiresAt: time.Now().Add(-time.Hour)}

	_, appErr := svc.Refresh(context.Background(), RefreshRequest{RefreshToken: raw})
	if appErr == nil || appErr.Code != errs.NewUnauthorized("").Code {
		t.Fatalf("expected unauthorized for expired token, got %v", appErr)
	}
}

func TestService_GetMe_NotFound(t *testing.T) {
	svc := testService(newMockUserRepo())
	_, appErr := svc.GetMe(context.Background(), "missing")
	if appErr == nil || appErr.Code != errs.NewNotFound("").Code {
		t.Fatalf("expected not found, got %v", appErr)
	}
}
