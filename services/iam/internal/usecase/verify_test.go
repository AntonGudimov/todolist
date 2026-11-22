package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthServiceSuite) TestVerifyOk() {
	createUserCommand := CreateUserCommand{
		Name:     "test name",
		Email:    "testEmail@yandex.ru",
		Password: "test password",
	}
	registerResult, registerErr := s.service.Register(s.ctx, createUserCommand)
	s.Require().NoError(registerErr)

	verifyErr := s.service.Verify(s.ctx, registerResult.Token)
	s.Require().NoError(verifyErr)

	user := s.uow.Db.Users[registerResult.UserID]

	emailVerification := s.uow.Db.EmailVerifications[registerResult.EmailVerificationID]

	s.Assert().WithinDuration(
		*emailVerification.UsedAt,
		time.Now().UTC(),
		time.Second,
	)
	s.Assert().Equal(true, user.IsVerified)
}

func (s *AuthServiceSuite) TestVerifyErr() {
	utcNow := time.Now().UTC()
	expiresAt := utcNow.Add(-1 * s.config.EmailVerificationTTL)
	usedAt := utcNow.Add(s.config.EmailVerificationTTL / 2)
	token, tokenErr := generateToken()
	s.Require().NoError(tokenErr)
	hashedToken := hashToken(token)

	tests := []struct {
		name              string
		rawToken          string
		emailVerification domain.EmailVerification
		expectedErr       error
	}{
		{
			name:     "email verification used",
			rawToken: token,
			emailVerification: domain.EmailVerification{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				Token:     hashedToken,
				ExpiresAt: utcNow.Add(s.config.EmailVerificationTTL),
				UsedAt:    &usedAt,
			},
			expectedErr: domain.ErrEmailVerificationUsed,
		},
		{
			name:     "email verification expired",
			rawToken: token,
			emailVerification: domain.EmailVerification{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				Token:     hashedToken,
				ExpiresAt: expiresAt,
				UsedAt:    nil,
			},
			expectedErr: domain.ErrEmailVerificationExpired,
		},
	}
	for _, tt := range tests {
		s.Run(
			tt.name,
			func() {
				s.uow.Db.EmailVerifications[tt.emailVerification.ID] = &tt.emailVerification
				err := s.service.Verify(s.ctx, tt.rawToken)
				s.Assert().ErrorIs(err, tt.expectedErr)

				s.uow.Db.Clear()
			},
		)
	}
}
