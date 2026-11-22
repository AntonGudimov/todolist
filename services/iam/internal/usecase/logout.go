package usecase

import "context"

type LogoutCommand struct {
	RefreshToken string
}

func (s *Service) Logout(ctx context.Context, cmd LogoutCommand) error {
	hashedRefreshToken := hashToken(cmd.RefreshToken)
	refreshTokenStored, refreshTokenErr := s.uow.Auth().GetRefreshTokenByHash(ctx, hashedRefreshToken)
	if refreshTokenErr != nil {
		return refreshTokenErr
	}
	revokeErr := s.uow.Auth().RevokeRefreshTokensByFamilyID(ctx, refreshTokenStored.FamilyID)
	return revokeErr
}
