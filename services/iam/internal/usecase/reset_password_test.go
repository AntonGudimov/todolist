package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthServiceSuite) ResetPasswordOk() {
	password := "test password"
	passwordHash, passwordHashErr := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	s.Require().NoError(passwordHashErr)
	updatedAt := time.Now().UTC().Add(-1 * time.Minute)
	user := domain.User{
		ID:           uuid.New(),
		Name:         "test name",
		Email:        "test email",
		PasswordHash: string(passwordHash),
		IsVerified:   true,
		CreatedAt:    time.Now().UTC().Add(-1 * time.Hour),
		UpdatedAt:    &updatedAt,
	}
	s.uow.Db.Users[user.ID] = &user

	loginResult, loginErr := s.service.Login(s.ctx, LoginUserCommand{Email: user.Email, Password: password})
	s.Require().NoError(loginErr)

	forgotPasswordResult, forgotPasswordErr := s.service.ForgotPassword(
		s.ctx,
		ForgotPasswordCommand{Email: user.Email},
	)
	s.Require().NoError(forgotPasswordErr)

	newPassword := "new password"
	resetPasswordErr := s.service.ResetPassword(
		s.ctx, ResetPasswordCommand{Token: forgotPasswordResult.Token, NewPassword: newPassword},
	)
	s.Require().NoError(resetPasswordErr)

	updatedUser := s.uow.Db.Users[user.ID]

	refreshToken, refreshTokenErr := s.uow.Auth().GetRefreshTokenByHash(s.ctx, hashToken(loginResult.RefreshToken))
	s.Require().NoError(refreshTokenErr)

	passwordResetToken, passwordResetTokenErr := s.uow.Auth().GetPasswordResetTokensByTokenHash(s.ctx, hashToken(forgotPasswordResult.Token))
	s.Require().NoError(passwordResetTokenErr)

	newPasswordHashErr := bcrypt.CompareHashAndPassword(
		[]byte(updatedUser.PasswordHash),
		[]byte(newPassword),
	)
	s.Assert().NoError(newPasswordHashErr)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*passwordResetToken.UsedAt,
		time.Second,
	)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
}

func (s *AuthServiceSuite) ResetPasswordTokenNotFoundErr() {
	err := s.service.ResetPassword(s.ctx, ResetPasswordCommand{Token: "test token", NewPassword: "test new password"})
	s.Assert().ErrorIs(err, domain.ErrPasswordResetTokenNotFound)
}

func (s *AuthServiceSuite) ResetPasswordTokenRevokedErr() {
	revokedAt := time.Now().UTC()
	token := "test token"
	passwordResetToken := domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().UTC().Add(s.config.PasswordResetTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: &revokedAt,
		UsedAt:    nil,
	}
	s.uow.Db.PasswordResetTokens[passwordResetToken.ID] = &passwordResetToken

	err := s.service.ResetPassword(s.ctx, ResetPasswordCommand{Token: token, NewPassword: "test new password"})
	s.Assert().ErrorIs(err, domain.ErrPasswordResetTokenRevoked)
}

func (s *AuthServiceSuite) ResetPasswordTokenUsedErr() {
	usedAt := time.Now().UTC()
	token := "test token"
	passwordResetToken := domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().UTC().Add(s.config.PasswordResetTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: nil,
		UsedAt:    &usedAt,
	}
	s.uow.Db.PasswordResetTokens[passwordResetToken.ID] = &passwordResetToken

	err := s.service.ResetPassword(s.ctx, ResetPasswordCommand{Token: token, NewPassword: "test new password"})
	s.Assert().ErrorIs(err, domain.ErrPasswordResetTokenUsed)
}

func (s *AuthServiceSuite) ResetPasswordTokenExpiredErr() {
	token := "test token"
	passwordResetToken := domain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		CreatedAt: time.Now().UTC(),
		RevokedAt: nil,
		UsedAt:    nil,
	}
	s.uow.Db.PasswordResetTokens[passwordResetToken.ID] = &passwordResetToken

	err := s.service.ResetPassword(s.ctx, ResetPasswordCommand{Token: token, NewPassword: "test new password"})
	s.Assert().ErrorIs(err, domain.ErrPasswordResetTokenExpired)
}
