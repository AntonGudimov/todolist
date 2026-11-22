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

func (a *AuthPostgresRepository) CreateEmailVerification(ctx context.Context, arg domain.CreateEmailVerificationParams) (domain.EmailVerification, error) {
	result, err := a.q.AddEmailVerification(
		ctx,
		generated.AddEmailVerificationParams{
			ID:        pgtype.UUID{Bytes: arg.ID, Valid: true},
			UserID:    pgtype.UUID{Bytes: arg.UserID, Valid: true},
			Token:     arg.Token,
			ExpiresAt: pgtype.Timestamptz{Time: arg.ExpiresAt, Valid: true},
			UsedAt:    pgtype.Timestamptz{Valid: false},
		},
	)
	if err != nil {
		return domain.EmailVerification{}, err
	}
	emailVerification := domain.EmailVerification{
		ID:        result.ID.Bytes,
		UserID:    result.UserID.Bytes,
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt.Time,
		UsedAt:    nil,
	}
	return emailVerification, nil
}

func (a *AuthPostgresRepository) GetEmailVerificationByToken(ctx context.Context, token string) (domain.EmailVerification, error) {
	result, err := a.q.GetEmailVerificationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EmailVerification{}, domain.ErrEmailVerificationNotFound
		}
		return domain.EmailVerification{}, err
	}
	var usedAt *time.Time
	if result.UsedAt.Valid {
		usedAt = &result.UsedAt.Time
	}
	return domain.EmailVerification{
		ID:        result.ID.Bytes,
		UserID:    result.UserID.Bytes,
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt.Time.UTC(),
		UsedAt:    usedAt,
	}, nil
}

func (a *AuthPostgresRepository) UseEmailVerification(ctx context.Context, arg domain.UseEmailVerificationParams) error {
	err := a.q.UseEmailVerification(
		ctx,
		generated.UseEmailVerificationParams{
			ID:     pgtype.UUID{Bytes: arg.ID, Valid: true},
			UsedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		},
	)
	return err
}

func (a *AuthPostgresRepository) CleanupEmailVerifications(ctx context.Context, cutoff time.Time, batchSize int) ([]uuid.UUID, error) {
	result, err := a.q.CleanupEmailVerifications(
		ctx,
		generated.CleanupEmailVerificationsParams{
			ExpiresAt: pgtype.Timestamptz{Time: cutoff.UTC(), Valid: true},
			Limit:     int32(batchSize),
		},
	)
	ids := make([]uuid.UUID, 0, cap(result))
	for _, id := range result {
		ids = append(ids, uuid.UUID(id.Bytes))
	}
	return ids, err
}
