package usecase

import (
	"context"
	"fmt"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

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

func (s *Service) Verify(
	ctx context.Context,
	token string,
) error {
	emailVerification, emailVerificationErr := s.uow.Auth().GetEmailVerificationByToken(
		ctx, hashToken(token),
	)
	if emailVerificationErr != nil {
		return emailVerificationErr
	}
	if emailVerification.UsedAt != nil {
		return fmt.Errorf("verify email: %v: %w", emailVerification.ID, domain.ErrEmailVerificationUsed)
	}
	usedAt := time.Now().UTC()
	if emailVerification.ExpiresAt.Before(usedAt) {
		return fmt.Errorf("verify email: %v: %w", emailVerification.ID, domain.ErrEmailVerificationExpired)
	}
	return s.uow.Do(
		ctx,
		func(uow domain.UnitOfWork) error {
			return s.verifyTx(ctx, emailVerification, usedAt)
		},
	)
}

func (s *Service) verifyTx(ctx context.Context, emailVerification domain.EmailVerification, usedAt time.Time) error {
	useEmailErr := s.uow.Auth().UseEmailVerification(
		ctx,
		domain.UseEmailVerificationParams{
			ID:     emailVerification.ID,
			UsedAt: usedAt,
		},
	)
	if useEmailErr != nil {
		return useEmailErr
	}
	verifyUserErr := s.uow.Auth().VerifyUser(
		ctx,
		emailVerification.UserID,
	)
	if verifyUserErr != nil {
		return verifyUserErr
	}
	return nil
}
