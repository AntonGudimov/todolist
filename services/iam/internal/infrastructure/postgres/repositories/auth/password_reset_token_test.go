package auth

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthPostgresRepositorySuite) TestAddPasswordResetTokenOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	_, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	tokenHash := "test token hash"
	createPasswordResetTokenParams := domain.CreatePasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	token, err := s.repo.CreatePasswordResetToken(s.ctx, createPasswordResetTokenParams)
	s.Require().NoError(err)

	s.Assert().Equal(createPasswordResetTokenParams.ID, token.ID)
	s.Assert().Equal(createPasswordResetTokenParams.UserID, token.UserID)
	s.Assert().Equal(createPasswordResetTokenParams.TokenHash, token.TokenHash)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		token.CreatedAt,
		time.Second,
	)
	s.Assert().WithinDuration(
		createPasswordResetTokenParams.ExpiresAt,
		token.ExpiresAt,
		time.Millisecond,
	)
}

func (s *AuthPostgresRepositorySuite) TestRevokePasswordResetTokensByUserIDOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	_, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	tokenHash := "test token hash"
	createPasswordResetTokenParams := domain.CreatePasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	token, err := s.repo.CreatePasswordResetToken(s.ctx, createPasswordResetTokenParams)
	s.Require().NoError(err)

	s.Assert().Equal(createPasswordResetTokenParams.ID, token.ID)
	s.Assert().Equal(createPasswordResetTokenParams.UserID, token.UserID)
	s.Assert().Equal(createPasswordResetTokenParams.TokenHash, token.TokenHash)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		token.CreatedAt,
		time.Second,
	)
	s.Assert().WithinDuration(
		createPasswordResetTokenParams.ExpiresAt,
		token.ExpiresAt,
		time.Millisecond,
	)
}

func (s *AuthPostgresRepositorySuite) TestGetPasswordResetTokensByTokenHash() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	_, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	tokenHash := "test token hash"
	createPasswordResetTokenParams := domain.CreatePasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	token, tokenErr := s.repo.CreatePasswordResetToken(s.ctx, createPasswordResetTokenParams)
	s.Require().NoError(tokenErr)

	passwordResetToken, err := s.repo.GetPasswordResetTokensByTokenHash(s.ctx, token.TokenHash)
	s.Require().NoError(err)

	s.Assert().Equal(token.ID, passwordResetToken.ID)
	s.Assert().Equal(token.UserID, passwordResetToken.UserID)
	s.Assert().Equal(token.TokenHash, passwordResetToken.TokenHash)
	s.Assert().WithinDuration(
		token.ExpiresAt,
		passwordResetToken.ExpiresAt,
		time.Millisecond,
	)
	s.Assert().WithinDuration(
		token.CreatedAt,
		passwordResetToken.CreatedAt,
		time.Millisecond,
	)
	s.Assert().Equal(token.UsedAt, passwordResetToken.UsedAt)
}

func (s *AuthPostgresRepositorySuite) TestGetPasswordResetTokensByTokenHashNotFoundErr() {
	_, err := s.repo.GetPasswordResetTokensByTokenHash(s.ctx, "test token hash")
	s.Assert().ErrorIs(err, domain.ErrPasswordResetTokenNotFound)
}

func (s *AuthPostgresRepositorySuite) TestUsePasswordResetTokenOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	_, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	tokenHash := "test token hash"
	createPasswordResetTokenParams := domain.CreatePasswordResetTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	token, tokenErr := s.repo.CreatePasswordResetToken(s.ctx, createPasswordResetTokenParams)
	s.Require().NoError(tokenErr)

	err := s.repo.UsePasswordResetToken(s.ctx, token.ID)
	s.Require().NoError(err)

	token, tokenErr = s.repo.GetPasswordResetTokensByTokenHash(s.ctx, token.TokenHash)
	s.Require().NoError(tokenErr)

	actualUsedAt := *token.UsedAt
	s.Assert().WithinDuration(time.Now().UTC(), actualUsedAt.UTC(), time.Second)
}
