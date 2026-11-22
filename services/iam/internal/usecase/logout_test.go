package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthServiceSuite) TestLogoutOk() {
	rawRefreshToken, rawRefreshTokenErr := generateRefreshToken()
	s.Require().NoError(rawRefreshTokenErr)

	revokedAt := time.Now().UTC()
	refreshToken := domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		FamilyID:  uuid.New(),
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: time.Now().UTC().Add(-1 * s.config.RefreshTokenTTL),
		CreatedAt: time.Now().UTC(),
		RevokedAt: &revokedAt,
	}
	s.uow.Db.RefreshTokens[refreshToken.ID] = &refreshToken

	err := s.service.Logout(s.ctx, LogoutCommand{RefreshToken: rawRefreshToken})
	s.Require().NoError(err)

	s.Assert().NotNil(refreshToken.RevokedAt)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*refreshToken.RevokedAt,
		time.Second,
	)
}
