package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthServiceSuite) TestForgotPasswordOk() {
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

	result, err := s.service.ForgotPassword(s.ctx, ForgotPasswordCommand{Email: user.Email})
	s.Require().NoError(err)

	passwordResetToken, passwordResetTokenErr := s.uow.Auth().GetPasswordResetTokensByTokenHash(s.ctx, hashToken(result.Token))
	s.Require().NoError(passwordResetTokenErr)

	s.Assert().Equal(user.ID, passwordResetToken.UserID)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		passwordResetToken.CreatedAt,
		time.Second,
	)
	s.Assert().WithinDuration(
		time.Now().UTC().Add(s.config.PasswordResetTokenTTL),
		passwordResetToken.ExpiresAt,
		time.Second,
	)
	s.Assert().Nil(passwordResetToken.RevokedAt)
	s.Assert().Nil(passwordResetToken.UsedAt)
	passwordResetInfo := s.emailSender.SentPasswordResetTokens[0]
	s.Assert().Contains(passwordResetInfo.ResetUrl, result.Token)
}

func (s *AuthServiceSuite) TestForgotPasswordSeveralTimes() {
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

	result1, err1 := s.service.ForgotPassword(s.ctx, ForgotPasswordCommand{Email: user.Email})
	s.Require().NoError(err1)

	result2, err2 := s.service.ForgotPassword(s.ctx, ForgotPasswordCommand{Email: user.Email})
	s.Require().NoError(err2)

	passwordResetToken1, passwordResetTokenErr1 := s.uow.Auth().GetPasswordResetTokensByTokenHash(s.ctx, hashToken(result1.Token))
	s.Require().NoError(passwordResetTokenErr1)

	passwordResetToken2, passwordResetTokenErr2 := s.uow.Auth().GetPasswordResetTokensByTokenHash(s.ctx, hashToken(result2.Token))
	s.Require().NoError(passwordResetTokenErr2)

	s.Assert().NotNil(passwordResetToken1.RevokedAt)
	s.Assert().Nil(passwordResetToken2.RevokedAt)
}

func (s *AuthServiceSuite) TestForgotPasswordUserNotFoundErr() {
	_, err := s.service.ForgotPassword(
		s.ctx, ForgotPasswordCommand{Email: "test email"},
	)
	s.Assert().ErrorIs(err, domain.ErrUserNotFound)
}
