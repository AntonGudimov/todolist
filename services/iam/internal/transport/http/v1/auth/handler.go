package auth

import (
	"log/slog"
	"todolist/internal/usecase"
)

type AuthHandler struct {
	authService usecase.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService usecase.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}
