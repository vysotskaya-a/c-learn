package errs

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func NewValidation(msg string, details any) *AppError {
	return &AppError{Code: http.StatusBadRequest, Error: "validation_error", Message: msg, Details: details}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Error: "unauthorized", Message: msg}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Error: "forbidden", Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Error: "not_found", Message: msg}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Error: "conflict", Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Error: "internal_error", Message: msg}
}
