package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
}

type CreateEmailVerificationParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type UseEmailVerificationParams struct {
	ID     uuid.UUID
	UsedAt time.Time
}

type CreateRefreshTokenParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FamilyID  uuid.UUID
	TokenHash string

	ExpiresAt time.Time
}

type CreatePasswordResetTokenParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

type AuthRepository interface {
	Register(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	CountUsersByEmail(ctx context.Context, email string) (int, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (User, error)
	CreateEmailVerification(ctx context.Context, arg CreateEmailVerificationParams) (EmailVerification, error)
	GetEmailVerificationByToken(ctx context.Context, token string) (EmailVerification, error)
	UseEmailVerification(ctx context.Context, arg UseEmailVerificationParams) error
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateRefreshTokenOk(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error
	RevokeRefreshTokensByFamilyID(ctx context.Context, familyID uuid.UUID) error
	CreatePasswordResetToken(ctx context.Context, params CreatePasswordResetTokenParams) (PasswordResetToken, error)
	RevokePasswordResetTokensByUserID(ctx context.Context, userID uuid.UUID) error
	GetPasswordResetTokensByTokenHash(ctx context.Context, tokenHash string) (PasswordResetToken, error)
	UpdatePassword(ctx context.Context, passwordHash string, userID uuid.UUID) error
	UsePasswordResetToken(ctx context.Context, ID uuid.UUID) error
	RevokeRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) error
	CleanupEmailVerifications(ctx context.Context, cutoff time.Time, batchSize int) ([]uuid.UUID, error)
}
