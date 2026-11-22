package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthServiceSuite) TestLoginOk() {
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

	loginUserResult, err := s.service.Login(s.ctx, LoginUserCommand{
		Email:    user.Email,
		Password: password,
	})
	s.Require().NoError(err)
	s.Assert().Len(s.uow.Db.RefreshTokens, 1)
	var storedRefreshToken domain.RefreshToken
	for _, refreshToken := range s.uow.Db.RefreshTokens {
		storedRefreshToken = *refreshToken
	}
	s.Assert().Equal(s.jwtService.AccessToken, loginUserResult.AccessToken)
	s.Assert().Equal(storedRefreshToken.TokenHash, hashToken(loginUserResult.RefreshToken))
}

func (s *AuthServiceSuite) TestLoginUserNotFoundErr() {
	token, err := s.service.Login(
		s.ctx,
		LoginUserCommand{
			Email:    "test email",
			Password: "test password",
		},
	)
	s.Assert().Equal(LoginUserResult{}, token)
	s.Assert().ErrorIs(err, domain.ErrInvalidCredentials)
}

func (s *AuthServiceSuite) TestLoginInvalidPasswordErr() {
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

	token, err := s.service.Login(s.ctx, LoginUserCommand{
		Email:    user.Email,
		Password: "invalid password",
	})
	s.Assert().Equal(LoginUserResult{}, token)
	s.Assert().ErrorIs(err, domain.ErrInvalidCredentials)
}
