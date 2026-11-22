package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FamilyID  uuid.UUID
	TokenHash string

	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
