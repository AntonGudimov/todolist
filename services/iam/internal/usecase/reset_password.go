package usecase

import (
	"context"
	"time"
	"todolist/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordCommand struct {
	Token       string
	NewPassword string
}

func (s *Service) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	passwordResetToken, passwordResetTokenErr := s.uow.Auth().GetPasswordResetTokensByTokenHash(
		ctx, hashToken(cmd.Token),
	)
	if passwordResetTokenErr != nil {
		return passwordResetTokenErr
	}

	if passwordResetToken.RevokedAt != nil {
		return domain.ErrPasswordResetTokenRevoked
	}
	if passwordResetToken.UsedAt != nil {
		return domain.ErrPasswordResetTokenUsed
	}
	if passwordResetToken.ExpiresAt.Before(time.Now().UTC()) {
		return domain.ErrPasswordResetTokenExpired
	}

	err := s.uow.Do(
		ctx,
		func(uow domain.UnitOfWork) error {
			err := s.resetPasswordTx(ctx, cmd, passwordResetToken, uow)
			return err
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) resetPasswordTx(ctx context.Context, cmd ResetPasswordCommand, passwordResetToken domain.PasswordResetToken, uow domain.UnitOfWork) error {
	passwordHash, passwordHashErr := bcrypt.GenerateFromPassword(
		[]byte(cmd.NewPassword),
		bcrypt.DefaultCost,
	)
	if passwordHashErr != nil {
		return passwordHashErr
	}
	updatePasswordErr := uow.Auth().UpdatePassword(ctx, string(passwordHash), passwordResetToken.UserID)
	if updatePasswordErr != nil {
		return updatePasswordErr
	}

	usePasswordResetTokenErr := uow.Auth().UsePasswordResetToken(ctx, passwordResetToken.ID)
	if usePasswordResetTokenErr != nil {
		return usePasswordResetTokenErr
	}

	revokeRefreshTokensErr := uow.Auth().RevokeRefreshTokenByUserID(ctx, passwordResetToken.UserID)
	if revokeRefreshTokensErr != nil {
		return revokeRefreshTokensErr
	}
	return nil
}
