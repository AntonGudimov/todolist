package auth

import (
	"errors"
	"net/http"
	"todolist/internal/domain"

	"github.com/gin-gonic/gin"
)

type VerifyUserRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) VerifyUser(c *gin.Context) {
	err := h.authService.Verify(
		c.Request.Context(),
		c.Query("token"),
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailVerificationNotFound):
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": err.Error()},
			)
		case errors.Is(err, domain.ErrEmailVerificationExpired) || errors.Is(err, domain.ErrEmailVerificationUsed):
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Something went wrong"},
			)
		}
	}
}
