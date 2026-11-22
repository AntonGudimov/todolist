package auth

import (
	"errors"
	"net/http"
	"time"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, refreshTokenErr := c.Cookie("refresh_token")
	if refreshTokenErr != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "refresh token required",
			},
		)
		return
	}
	refreshTokenResult, err := h.authService.Refresh(
		c.Request.Context(),
		usecase.RefreshTokenCommand{
			RefreshToken: refreshToken,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRefreshTokenNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
		case errors.Is(err, domain.ErrRefreshTokenInvalid):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}
	c.SetCookie(
		"refresh_token",
		refreshTokenResult.RefreshToken,
		int(time.Until(refreshTokenResult.RefreshTokenExpiresAt).Seconds()),
		"/auth",
		"",
		true,
		true,
	)
	c.JSON(
		http.StatusOK,
		gin.H{
			"accessToken": refreshTokenResult.AccessToken,
			"type":        "Bearer",
		},
	)
}
