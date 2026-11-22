package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	loginResult, loginResultErr := h.authService.Login(c.Request.Context(), usecase.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if loginResultErr != nil {
		h.logger.ErrorContext(
			c.Request.Context(),
			"login operation failed",
			slog.String("error", loginResultErr.Error()),
		)
		switch {
		case errors.Is(loginResultErr, domain.ErrInvalidCredentials):
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": loginResultErr.Error(),
				},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Something went wrong"},
			)
		}
		return
	}
	c.SetCookie(
		"refresh_token",
		loginResult.RefreshToken,
		int(time.Until(loginResult.RefreshTokenExpiresAt).Seconds()),
		"/api/v1/auth",
		"",
		true,
		true,
	)
	c.JSON(
		http.StatusOK,
		gin.H{
			"accessToken": loginResult.AccessToken,
			"type":        "Bearer",
		},
	)
}
