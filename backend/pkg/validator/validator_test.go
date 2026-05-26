package validator

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{name: "valid", email: "user@example.com", wantErr: false},
		{name: "empty", email: "", wantErr: true},
		{name: "invalid format", email: "not-an-email", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "valid 8 chars", password: "12345678", wantErr: false},
		{name: "valid unicode", password: "пароль12", wantErr: false},
		{name: "too short", password: "1234567", wantErr: true},
		{name: "empty", password: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{name: "valid min", username: "abc", wantErr: false},
		{name: "valid max", username: strings.Repeat("a", 100), wantErr: false},
		{name: "too short", username: "ab", wantErr: true},
		{name: "too long", username: strings.Repeat("a", 101), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSourceCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "valid", code: "int main() {}", wantErr: false},
		{name: "empty", code: "", wantErr: true},
		{name: "whitespace only", code: "   \n\t", wantErr: true},
		{name: "max size", code: strings.Repeat("a", 50*1024), wantErr: false},
		{name: "over max size", code: strings.Repeat("a", 50*1024+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSourceCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSourceCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	if err := ValidateRequired("value", "field"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := ValidateRequired("  ", "field"); err == nil {
		t.Fatal("expected error for blank value")
	}
}
