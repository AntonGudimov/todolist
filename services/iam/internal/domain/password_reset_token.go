package domain

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string

	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
	UsedAt    *time.Time
}
