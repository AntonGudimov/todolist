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

func (a *AuthPostgresRepository) CreateRefreshTokenOk(ctx context.Context, params domain.CreateRefreshTokenParams) (domain.RefreshToken, error) {
	refreshToken, err := a.q.AddRefreshToken(
		ctx,
		generated.AddRefreshTokenParams{
			ID:        pgtype.UUID{Bytes: params.ID, Valid: true},
			UserID:    pgtype.UUID{Bytes: params.UserID, Valid: true},
			FamilyID:  pgtype.UUID{Bytes: params.FamilyID, Valid: true},
			TokenHash: params.TokenHash,
			ExpiresAt: pgtype.Timestamptz{Time: params.ExpiresAt, Valid: true},
		},
	)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return domain.RefreshToken{
		ID:        refreshToken.ID.Bytes,
		UserID:    refreshToken.UserID.Bytes,
		FamilyID:  refreshToken.FamilyID.Bytes,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt.Time.UTC(),
		CreatedAt: refreshToken.CreatedAt.Time.UTC(),
		RevokedAt: nil,
	}, nil
}

func (a *AuthPostgresRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	refreshToken, err := a.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
		}
		return domain.RefreshToken{}, err
	}
	var revokedAt *time.Time
	if refreshToken.RevokedAt.Valid {
		revokedAt = &refreshToken.RevokedAt.Time
	}
	return domain.RefreshToken{
		ID:        refreshToken.ID.Bytes,
		UserID:    refreshToken.UserID.Bytes,
		FamilyID:  refreshToken.FamilyID.Bytes,
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt.Time,
		CreatedAt: refreshToken.CreatedAt.Time,
		RevokedAt: revokedAt,
	}, nil
}

func (a *AuthPostgresRepository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	err := a.q.RevokeRefreshTokenByHash(ctx, tokenHash)
	return err
}

func (a *AuthPostgresRepository) RevokeRefreshTokensByFamilyID(ctx context.Context, familyID uuid.UUID) error {
	err := a.q.RevokeRefreshTokensByFamilyId(ctx, pgtype.UUID{Bytes: familyID, Valid: true})
	return err
}

func (a *AuthPostgresRepository) RevokeRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) error {
	err := a.q.RevokeRefreshTokensByUserId(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	return err
}
