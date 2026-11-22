package auth

import (
	"errors"
	"net/http"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ForgotRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	result, err := h.authService.ForgotPassword(c.Request.Context(), usecase.ForgotPasswordCommand{Email: req.Email})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "user not found"},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Something went wrong"},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"token": result.Token,
		},
	)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	err := h.authService.ResetPassword(
		c.Request.Context(), usecase.ResetPasswordCommand{Token: req.Token, NewPassword: req.NewPassword},
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPasswordResetTokenNotFound):
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "password reset token not found"},
			)
		case errors.Is(err, domain.ErrPasswordResetTokenRevoked):
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "password reset token is revoked"},
			)
		case errors.Is(err, domain.ErrPasswordResetTokenUsed):
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "password reset token already used"},
			)
		case errors.Is(err, domain.ErrPasswordResetTokenExpired):
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "password reset token expired"},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Something went wrong"},
			)
		}
		return
	}
	c.JSON(
		http.StatusNoContent,
		gin.H{},
	)
}
