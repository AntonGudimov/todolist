package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

type ForgotPasswordCommand struct {
	Email string
}

type ForgotPasswordResult struct {
	Token string
}

func (s *Service) ForgotPassword(ctx context.Context, cmd ForgotPasswordCommand) (ForgotPasswordResult, error) {
	user, userErr := s.uow.Auth().GetUserByEmail(ctx, cmd.Email)
	if userErr != nil {
		return ForgotPasswordResult{}, userErr
	}

	var forgotPasswordResult ForgotPasswordResult
	err := s.uow.Do(
		ctx,
		func(uow domain.UnitOfWork) error {
			result, err := s.forgotPasswordTx(ctx, user)
			if err != nil {
				return err
			}
			forgotPasswordResult = result
			return nil
		},
	)
	if err != nil {
		return ForgotPasswordResult{}, err
	}
	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.config.AuthPublicUrl,
		url.QueryEscape(forgotPasswordResult.Token),
	)
	s.emailSender.SendForgotPasswordEmail(
		ctx,
		user.Email,
		int(s.config.PasswordResetTokenTTL.Minutes()),
		resetURL,
	)
	return forgotPasswordResult, nil
}

func (s *Service) forgotPasswordTx(ctx context.Context, user domain.User) (ForgotPasswordResult, error) {
	revokeErr := s.uow.Auth().RevokePasswordResetTokensByUserID(ctx, user.ID)
	if revokeErr != nil {
		return ForgotPasswordResult{}, revokeErr
	}
	passwordResetToken, passwordResetTokenErr := generateToken()
	if passwordResetTokenErr != nil {
		return ForgotPasswordResult{}, passwordResetTokenErr
	}
	_, err := s.uow.Auth().CreatePasswordResetToken(ctx, domain.CreatePasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashToken(passwordResetToken),
		ExpiresAt: time.Now().UTC().Add(s.config.PasswordResetTokenTTL),
	})
	if err != nil {
		return ForgotPasswordResult{}, err
	}
	return ForgotPasswordResult{Token: passwordResetToken}, nil
}
