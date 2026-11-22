package domain

import "github.com/google/uuid"

type JWTService interface {
	CreateAccessToken(userID uuid.UUID) (string, error)
}
