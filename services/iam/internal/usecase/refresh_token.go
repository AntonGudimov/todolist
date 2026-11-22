package usecase

import (
	"context"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

type RefreshTokenCommand struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (s *Service) Refresh(ctx context.Context, cmd RefreshTokenCommand) (RefreshTokenResult, error) {
	hashedRefreshToken := hashToken(cmd.RefreshToken)
	refreshTokenStored, refreshTokenErr := s.uow.Auth().GetRefreshTokenByHash(ctx, hashedRefreshToken)
	if refreshTokenErr != nil {
		return RefreshTokenResult{}, refreshTokenErr
	}
	if refreshTokenStored.RevokedAt != nil {
		revokeErr := s.uow.Auth().RevokeRefreshTokensByFamilyID(ctx, refreshTokenStored.FamilyID)
		if revokeErr != nil {
			return RefreshTokenResult{}, revokeErr
		}
		return RefreshTokenResult{}, domain.ErrRefreshTokenInvalid
	}
	if refreshTokenStored.ExpiresAt.Before(time.Now().UTC()) {
		return RefreshTokenResult{}, domain.ErrRefreshTokenInvalid
	}
	var refrestTokenResult RefreshTokenResult
	err := s.uow.Do(
		ctx,
		func(uow domain.UnitOfWork) error {
			accessToken, refreshToken, err := s.refreshTx(ctx, refreshTokenStored)
			refrestTokenResult = RefreshTokenResult{
				AccessToken:           accessToken,
				RefreshToken:          refreshToken,
				RefreshTokenExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
			}
			return err
		},
	)
	if err != nil {
		return RefreshTokenResult{}, err
	}
	return refrestTokenResult, nil
}

func (s *Service) refreshTx(ctx context.Context, refreshTokenStored domain.RefreshToken) (string, string, error) {
	revokeErr := s.uow.Auth().RevokeRefreshTokenByHash(ctx, refreshTokenStored.TokenHash)
	if revokeErr != nil {
		return "", "", revokeErr
	}
	refreshToken, refreshTokenErr := generateRefreshToken()
	if refreshTokenErr != nil {
		return "", "", refreshTokenErr
	}
	accessToken, accessTokenErr := s.jwtService.CreateAccessToken(refreshTokenStored.UserID)
	if accessTokenErr != nil {
		return "", "", accessTokenErr
	}
	_, err := s.uow.Auth().CreateRefreshTokenOk(
		ctx,
		domain.CreateRefreshTokenParams{
			ID:        uuid.New(),
			UserID:    refreshTokenStored.UserID,
			FamilyID:  refreshTokenStored.FamilyID,
			TokenHash: hashToken(refreshToken),
			ExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
		},
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
