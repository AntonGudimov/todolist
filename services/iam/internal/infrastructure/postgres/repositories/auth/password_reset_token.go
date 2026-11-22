package auth

import (
	"context"
	"errors"
	"time"
	"todolist/internal/domain"
	"todolist/internal/infrastructure/postgres/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a *AuthPostgresRepository) CreatePasswordResetToken(ctx context.Context, params domain.CreatePasswordResetTokenParams) (domain.PasswordResetToken, error) {
	token, err := a.q.AddPasswordResetToken(ctx, generated.AddPasswordResetTokenParams{
		ID:        pgtype.UUID{Bytes: params.ID, Valid: true},
		UserID:    pgtype.UUID{Bytes: params.UserID, Valid: true},
		TokenHash: params.TokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: params.ExpiresAt, Valid: true},
	})
	if err != nil {
		return domain.PasswordResetToken{}, err
	}
	return domain.PasswordResetToken{
		ID:        token.ID.Bytes,
		UserID:    token.UserID.Bytes,
		TokenHash: token.TokenHash,
		CreatedAt: token.CreatedAt.Time.UTC(),
		ExpiresAt: token.ExpiresAt.Time.UTC(),
		UsedAt:    nil,
	}, nil
}

func (a *AuthPostgresRepository) RevokePasswordResetTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	err := a.q.RevokePasswordResetTokensByUserId(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	return err
}

func (a *AuthPostgresRepository) GetPasswordResetTokensByTokenHash(ctx context.Context, tokenHash string) (domain.PasswordResetToken, error) {
	passwordResetToken, err := a.q.GetPasswordResetTokenByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PasswordResetToken{}, domain.ErrPasswordResetTokenNotFound
		}
		return domain.PasswordResetToken{}, err
	}
	var revokedAt *time.Time
	if passwordResetToken.RevokedAt.Valid {
		revokedAt = &passwordResetToken.RevokedAt.Time
	}
	var usedAt *time.Time
	if passwordResetToken.UsedAt.Valid {
		usedAt = &passwordResetToken.UsedAt.Time
	}

	return domain.PasswordResetToken{
		ID:        passwordResetToken.ID.Bytes,
		UserID:    passwordResetToken.UserID.Bytes,
		TokenHash: passwordResetToken.TokenHash,
		ExpiresAt: passwordResetToken.ExpiresAt.Time,
		CreatedAt: passwordResetToken.CreatedAt.Time,
		RevokedAt: revokedAt,
		UsedAt:    usedAt,
	}, nil
}

func (a *AuthPostgresRepository) UsePasswordResetToken(ctx context.Context, ID uuid.UUID) error {
	err := a.q.UsePasswordResetToken(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	return err
}
