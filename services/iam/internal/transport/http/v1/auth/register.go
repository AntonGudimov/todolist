package auth

import (
	"errors"
	"net/http"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12,max=128"`
}

type RegisterResponse struct {
	ID uuid.UUID `json:"ID" binding:"required,uuid"`
}

func (h *AuthHandler) CreateUser(
	c *gin.Context,
) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	registerResult, err := h.authService.Register(
		c.Request.Context(),
		usecase.CreateUserCommand{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(
				http.StatusConflict,
				gin.H{
					"message": err.Error(),
				},
			)
		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "internal server error",
				},
			)
		}
		return
	}

	c.JSON(
		http.StatusCreated,
		RegisterResponse{
			ID: registerResult.UserID,
		},
	)
}
