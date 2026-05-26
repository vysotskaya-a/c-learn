package errs

import (
	"net/http"
	"testing"
)

func TestAppErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		code int
	}{
		{"validation", NewValidation("bad", nil), http.StatusBadRequest},
		{"unauthorized", NewUnauthorized("no auth"), http.StatusUnauthorized},
		{"forbidden", NewForbidden("denied"), http.StatusForbidden},
		{"not found", NewNotFound("missing"), http.StatusNotFound},
		{"conflict", NewConflict("exists"), http.StatusConflict},
		{"internal", NewInternal("fail"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Fatalf("got code %d, want %d", tt.err.Code, tt.code)
			}
			if tt.err.Error == "" || tt.err.Message == "" {
				t.Fatal("error and message must be set")
			}
		})
	}
}
