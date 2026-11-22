package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthServiceSuite) TestRefreshOk() {
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

	rawRefreshToken, rawRefreshTokenErr := generateRefreshToken()
	s.Require().NoError(rawRefreshTokenErr)

	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		FamilyID:  uuid.New(),
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: nil,
	}
	s.uow.Db.RefreshTokens[refreshToken.ID] = &refreshToken

	result, err := s.service.Refresh(s.ctx, RefreshTokenCommand{
		RefreshToken: rawRefreshToken,
	})
	newRefreshToken, _ := s.uow.Auth().GetRefreshTokenByHash(
		s.ctx, hashToken(result.RefreshToken),
	)
	s.Require().NoError(err)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
	s.Assert().Equal(refreshToken.UserID, newRefreshToken.UserID)
	s.Assert().Equal(refreshToken.FamilyID, newRefreshToken.FamilyID)
}

func (s *AuthServiceSuite) TestRefreshTokenNotFoundErr() {
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

	_, err := s.service.Refresh(s.ctx, RefreshTokenCommand{
		RefreshToken: "test refresh token",
	})
	s.Assert().ErrorIs(err, domain.ErrRefreshTokenNotFound)
}

func (s *AuthServiceSuite) TestRefreshTokenRevokedErr() {
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

	rawRefreshToken, rawRefreshTokenErr := generateRefreshToken()
	s.Require().NoError(rawRefreshTokenErr)

	revokedAt := time.Now().UTC()
	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		FamilyID:  uuid.New(),
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: &revokedAt,
	}
	s.uow.Db.RefreshTokens[refreshToken.ID] = &refreshToken

	newRefreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    refreshToken.UserID,
		FamilyID:  refreshToken.FamilyID,
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: &revokedAt,
	}
	s.uow.Db.RefreshTokens[newRefreshToken.ID] = &refreshToken

	_, err := s.service.Refresh(s.ctx, RefreshTokenCommand{
		RefreshToken: rawRefreshToken,
	})
	s.Assert().ErrorIs(err, domain.ErrRefreshTokenInvalid)
	s.Assert().NotNil(newRefreshToken.RevokedAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*newRefreshToken.RevokedAt,
		time.Second,
	)
}

func (s *AuthServiceSuite) TestRefreshTokenExpiredErr() {
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

	rawRefreshToken, rawRefreshTokenErr := generateRefreshToken()
	s.Require().NoError(rawRefreshTokenErr)

	revokedAt := time.Now().UTC()
	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		FamilyID:  uuid.New(),
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(-1 * s.config.RefreshTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: &revokedAt,
	}
	s.uow.Db.RefreshTokens[refreshToken.ID] = &refreshToken

	_, err := s.service.Refresh(s.ctx, RefreshTokenCommand{
		RefreshToken: rawRefreshToken,
	})
	s.Assert().ErrorIs(err, domain.ErrRefreshTokenInvalid)
}
