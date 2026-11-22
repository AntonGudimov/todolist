package usecase

import (
	"context"
	"errors"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserCommand struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (s *Service) Login(ctx context.Context, cmd LoginUserCommand) (LoginUserResult, error) {
	user, userErr := s.uow.Auth().GetUserByEmail(ctx, cmd.Email)
	if userErr != nil {
		if errors.Is(userErr, domain.ErrUserNotFound) {
			return LoginUserResult{}, domain.ErrInvalidCredentials
		}
		return LoginUserResult{}, userErr
	}
	passwordHashErr := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(cmd.Password),
	)
	if passwordHashErr != nil {
		return LoginUserResult{}, domain.ErrInvalidCredentials
	}

	accessToken, accessTokenErr := s.jwtService.CreateAccessToken(user.ID)
	if accessTokenErr != nil {
		return LoginUserResult{}, accessTokenErr
	}

	generatedRefreshToken, generatedRefreshTokenErr := generateRefreshToken()
	if generatedRefreshTokenErr != nil {
		return LoginUserResult{}, generatedRefreshTokenErr
	}
	createRefreshTokenParams := domain.CreateRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FamilyID:  uuid.New(),
		TokenHash: hashToken(generatedRefreshToken),
		ExpiresAt: time.Now().UTC().Add(s.config.RefreshTokenTTL),
	}
	_, refreshTokenErr := s.uow.Auth().CreateRefreshTokenOk(
		ctx,
		createRefreshTokenParams,
	)
	if refreshTokenErr != nil {
		return LoginUserResult{}, refreshTokenErr
	}
	return LoginUserResult{
		AccessToken:           accessToken,
		RefreshToken:          generatedRefreshToken,
		RefreshTokenExpiresAt: createRefreshTokenParams.ExpiresAt,
	}, nil
}
