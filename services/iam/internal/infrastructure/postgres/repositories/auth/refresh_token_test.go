package auth

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthPostgresRepositorySuite) TestCreateRefreshToken() {
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

	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		FamilyID:  uuid.New(),
		TokenHash: "test token hash",
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
	refreshToken, err := s.repo.CreateRefreshTokenOk(s.ctx, createRefreshTokenParams)
	s.Require().NoError(err)

	s.Assert().Equal(createRefreshTokenParams.ID, refreshToken.ID)
	s.Assert().Equal(createRefreshTokenParams.UserID, refreshToken.UserID)
	s.Assert().Equal(createRefreshTokenParams.TokenHash, refreshToken.TokenHash)
	s.Assert().Equal(createRefreshTokenParams.FamilyID, refreshToken.FamilyID)
	s.Assert().Equal(createRefreshTokenParams.ExpiresAt, createRefreshTokenParams.ExpiresAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		refreshToken.CreatedAt,
		time.Second,
	)
	s.Assert().Nil(refreshToken.RevokedAt)
}

func (s *AuthPostgresRepositorySuite) TestGetRefreshTokenByHashOk() {
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

	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		FamilyID:  uuid.New(),
		TokenHash: "test token hash",
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
	_, err := s.repo.CreateRefreshTokenOk(s.ctx, createRefreshTokenParams)
	s.Require().NoError(err)

	refreshToken, refreshTokenErr := s.repo.GetRefreshTokenByHash(s.ctx, createRefreshTokenParams.TokenHash)
	s.Require().NoError(refreshTokenErr)

	s.Assert().Equal(createRefreshTokenParams.ID, refreshToken.ID)
}

func (s *AuthPostgresRepositorySuite) TestGetRefreshTokenByHashNotFound() {
	_, refreshTokenErr := s.repo.GetRefreshTokenByHash(s.ctx, "test hash")
	s.Assert().ErrorIs(refreshTokenErr, domain.ErrRefreshTokenNotFound)
}

func (s *AuthPostgresRepositorySuite) TestRevokeRefreshTokenByHashOk() {
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

	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		FamilyID:  uuid.New(),
		TokenHash: "test token hash",
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
	_, createTokenErr := s.repo.CreateRefreshTokenOk(s.ctx, createRefreshTokenParams)
	s.Require().NoError(createTokenErr)

	err := s.repo.RevokeRefreshTokenByHash(s.ctx, createRefreshTokenParams.TokenHash)
	s.Require().NoError(err)

	refreshToken, refreshTokenErr := s.repo.GetRefreshTokenByHash(s.ctx, createRefreshTokenParams.TokenHash)
	s.Require().NoError(refreshTokenErr)

	s.Assert().Equal(createRefreshTokenParams.ID, refreshToken.ID)
	s.Assert().NotNil(refreshToken.RevokedAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
}

func (s *AuthPostgresRepositorySuite) TestRevokeRefreshTokensByFamilyIDOk() {
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

	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		FamilyID:  uuid.New(),
		TokenHash: "test token hash",
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
	_, createTokenErr := s.repo.CreateRefreshTokenOk(s.ctx, createRefreshTokenParams)
	s.Require().NoError(createTokenErr)

	revokeErr := s.repo.RevokeRefreshTokensByFamilyID(s.ctx, createRefreshTokenParams.FamilyID)
	s.Require().NoError(revokeErr)

	refreshToken, refreshTokenErr := s.repo.GetRefreshTokenByHash(s.ctx, createRefreshTokenParams.TokenHash)
	s.Require().NoError(refreshTokenErr)

	s.Assert().Equal(createRefreshTokenParams.ID, refreshToken.ID)
	s.Assert().NotNil(refreshToken.RevokedAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
}

func (s *AuthPostgresRepositorySuite) TestRevokeRefreshTokensByUserIDOk() {
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

	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    createUserParams.ID,
		FamilyID:  uuid.New(),
		TokenHash: "test token hash",
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
	_, createTokenErr := s.repo.CreateRefreshTokenOk(s.ctx, createRefreshTokenParams)
	s.Require().NoError(createTokenErr)

	revokeErr := s.repo.RevokeRefreshTokenByUserID(s.ctx, createRefreshTokenParams.UserID)
	s.Require().NoError(revokeErr)

	refreshToken, refreshTokenErr := s.repo.GetRefreshTokenByHash(s.ctx, createRefreshTokenParams.TokenHash)
	s.Require().NoError(refreshTokenErr)

	s.Assert().Equal(createRefreshTokenParams.ID, refreshToken.ID)
	s.Assert().NotNil(refreshToken.RevokedAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
}
