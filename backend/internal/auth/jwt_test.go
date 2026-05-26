package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndValidateAccessToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")
	token, err := mgr.GenerateAccessToken("user-1", "student")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := mgr.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if claims.UserID != "user-1" || claims.Role != "student" {
		t.Fatalf("claims = %+v, want user-1/student", claims)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")
	if _, err := mgr.ValidateAccessToken("invalid.token.value"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	mgr1 := NewJWTManager("secret-a")
	mgr2 := NewJWTManager("secret-b")

	token, err := mgr1.GenerateAccessToken("user-1", "student")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if _, err := mgr2.ValidateAccessToken(token); err == nil {
		t.Fatal("expected error when validating with different secret")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	mgr := NewJWTManager("test-secret")
	raw, hash, expiresAt, err := mgr.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if len(raw) != 64 {
		t.Fatalf("raw token length = %d, want 64 hex chars", len(raw))
	}
	if hash != HashRefreshToken(raw) {
		t.Fatal("hash mismatch")
	}
	if expiresAt.Before(time.Now()) {
		t.Fatal("expiresAt must be in the future")
	}
}

func TestHashRefreshToken_Deterministic(t *testing.T) {
	raw := "abc123"
	if HashRefreshToken(raw) != HashRefreshToken(raw) {
		t.Fatal("hash must be deterministic")
	}
	if HashRefreshToken("a") == HashRefreshToken("b") {
		t.Fatal("different tokens must produce different hashes")
	}
}
