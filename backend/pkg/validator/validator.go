package validator

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

func ValidateUsername(username string) error {
	l := utf8.RuneCountInString(username)
	if l < 3 || l > 100 {
		return fmt.Errorf("username must be between 3 and 100 characters")
	}
	return nil
}

func ValidateSourceCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("source_code is required")
	}
	if len(code) > 50*1024 {
		return fmt.Errorf("source_code must be at most 50 KB")
	}
	return nil
}

func ValidateRequired(field, name string) error {
	if strings.TrimSpace(field) == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}
