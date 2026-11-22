package fakes

import (
	"context"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

type FakeAuthRepository struct {
	db *FakeDB
}

func NewFakeAuthRepository(db *FakeDB) domain.AuthRepository {
	return &FakeAuthRepository{
		db: db,
	}
}

func (a *FakeAuthRepository) Register(ctx context.Context, arg domain.CreateUserParams) (uuid.UUID, error) {
	user := domain.User{
		ID:           uuid.New(),
		Name:         arg.Name,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
		IsVerified:   false,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    nil,
	}
	a.db.Users[user.ID] = &user
	return user.ID, nil
}

func (a *FakeAuthRepository) CountUsersByEmail(ctx context.Context, email string) (int, error) {
	count := 0
	for _, user := range a.db.Users {
		if user.Email == email {
			count++
		}
	}
	return count, nil
}

func (a *FakeAuthRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	user := a.db.Users[ID]
	return *user, nil
}

func (a *FakeAuthRepository) CreateEmailVerification(ctx context.Context, arg domain.CreateEmailVerificationParams) (domain.EmailVerification, error) {
	emailVerification := domain.EmailVerification{
		ID:        arg.ID,
		UserID:    arg.UserID,
		Token:     arg.Token,
		ExpiresAt: arg.ExpiresAt,
		UsedAt:    nil,
	}
	a.db.EmailVerifications[arg.ID] = &emailVerification
	return emailVerification, nil
}

func (a *FakeAuthRepository) GetEmailVerificationByToken(ctx context.Context, token string) (domain.EmailVerification, error) {
	for _, emailVerification := range a.db.EmailVerifications {
		if emailVerification.Token == token {
			return *emailVerification, nil
		}
	}
	return domain.EmailVerification{}, domain.ErrEmailVerificationNotFound
}

func (a *FakeAuthRepository) UseEmailVerification(ctx context.Context, arg domain.UseEmailVerificationParams) error {
	emailVerification := a.db.EmailVerifications[arg.ID]
	emailVerification.UsedAt = &arg.UsedAt
	return nil
}

func (a *FakeAuthRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	user := a.db.Users[userID]
	user.IsVerified = true
	t := time.Now().UTC()
	user.UpdatedAt = &t
	return nil
}

func (a *FakeAuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	for _, user := range a.db.Users {
		if user.Email == email {
			return domain.User{
					ID:           user.ID,
					Name:         user.Name,
					Email:        user.Email,
					PasswordHash: user.PasswordHash,
					IsVerified:   user.IsVerified,
					CreatedAt:    user.CreatedAt,
					UpdatedAt:    user.UpdatedAt,
				},
				nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (a *FakeAuthRepository) CreateRefreshTokenOk(ctx context.Context, arg domain.CreateRefreshTokenParams) (domain.RefreshToken, error) {
	refreshToken := domain.RefreshToken{
		ID:        arg.ID,
		UserID:    arg.UserID,
		FamilyID:  arg.FamilyID,
		TokenHash: arg.TokenHash,
		ExpiresAt: arg.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		RevokedAt: nil,
	}
	a.db.RefreshTokens[refreshToken.ID] = &refreshToken
	return refreshToken, nil
}

func (a *FakeAuthRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	for _, refreshToken := range a.db.RefreshTokens {
		if refreshToken.TokenHash == tokenHash {
			return *refreshToken, nil
		}
	}
	return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
}

func (a *FakeAuthRepository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	for _, refreshToken := range a.db.RefreshTokens {
		if refreshToken.TokenHash == tokenHash {
			utcNow := time.Now().UTC()
			refreshToken.RevokedAt = &utcNow
			return nil
		}
	}
	return nil
}

func (a *FakeAuthRepository) RevokeRefreshTokensByFamilyID(ctx context.Context, familyID uuid.UUID) error {
	for _, refreshToken := range a.db.RefreshTokens {
		if refreshToken.FamilyID == familyID {
			utcNow := time.Now().UTC()
			refreshToken.RevokedAt = &utcNow
		}
	}
	return nil
}

func (a *FakeAuthRepository) CreatePasswordResetToken(ctx context.Context, params domain.CreatePasswordResetTokenParams) (domain.PasswordResetToken, error) {
	passwordResetToken := domain.PasswordResetToken{
		ID:        params.ID,
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: time.Now().UTC(),
		UsedAt:    nil,
	}
	a.db.PasswordResetTokens[passwordResetToken.ID] = &passwordResetToken
	return passwordResetToken, nil
}

func (a *FakeAuthRepository) RevokePasswordResetTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	for _, passwordResetToken := range a.db.PasswordResetTokens {
		if passwordResetToken.UserID == userID {
			utcNow := time.Now().UTC()
			passwordResetToken.RevokedAt = &utcNow
		}
	}
	return nil
}

func (a *FakeAuthRepository) GetPasswordResetTokensByTokenHash(ctx context.Context, tokenHash string) (domain.PasswordResetToken, error) {
	for _, passwordResetToken := range a.db.PasswordResetTokens {
		if passwordResetToken.TokenHash == tokenHash {
			return *passwordResetToken, nil
		}
	}
	return domain.PasswordResetToken{}, domain.ErrPasswordResetTokenNotFound
}

func (a *FakeAuthRepository) UpdatePassword(ctx context.Context, passwordHash string, userID uuid.UUID) error {
	for _, user := range a.db.Users {
		if user.ID == userID {
			user.PasswordHash = passwordHash
			return nil
		}
	}
	return nil
}

func (a *FakeAuthRepository) UsePasswordResetToken(ctx context.Context, ID uuid.UUID) error {
	for _, passwordResetToken := range a.db.PasswordResetTokens {
		if passwordResetToken.ID == ID {
			utcNow := time.Now().UTC()
			passwordResetToken.UsedAt = &utcNow
			return nil
		}
	}
	return nil
}

func (a *FakeAuthRepository) RevokeRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) error {
	for _, refreshToken := range a.db.RefreshTokens {
		if refreshToken.UserID == userID {
			utcNow := time.Now().UTC()
			refreshToken.RevokedAt = &utcNow
		}
	}
	return nil
}

func (a *FakeAuthRepository) CleanupEmailVerifications(ctx context.Context, cutoff time.Time, batchSize int) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(a.db.EmailVerifications))
	for _, emailVerification := range a.db.EmailVerifications {
		if len(result) == batchSize {
			return result, nil
		}
		if emailVerification.ExpiresAt.Equal(cutoff) || emailVerification.ExpiresAt.After(cutoff) {
			result = append(result, emailVerification.ID)
			delete(a.db.EmailVerifications, emailVerification.ID)
		}
	}
	return result, nil
}
