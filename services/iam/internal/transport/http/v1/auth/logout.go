package auth

import (
	"net/http"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshTokenCookie, refreshTokenErr := c.Cookie("refresh_token")
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/auth",
		"",
		true,
		true,
	)
	if refreshTokenErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
		return
	}

	err := h.authService.Logout(c.Request.Context(), usecase.LogoutCommand{RefreshToken: refreshTokenCookie})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
