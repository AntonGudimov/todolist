package auth

import (
	"context"
	"errors"
	"time"
	"todolist/internal/domain"
	"todolist/internal/infrastructure/postgres/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a *AuthPostgresRepository) Register(ctx context.Context, arg domain.CreateUserParams) (uuid.UUID, error) {
	user, err := a.q.AddUser(ctx, generated.AddUserParams{
		ID:         pgtype.UUID{Bytes: arg.ID, Valid: true},
		Name:       arg.Name,
		Email:      arg.Email,
		Password:   arg.PasswordHash,
		IsVerified: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID.Bytes, nil
}

func (a *AuthPostgresRepository) CountUsersByEmail(ctx context.Context, email string) (int, error) {
	userCount, errCount := a.q.CountUsersByEmail(ctx, email)
	if errCount != nil {
		return 0, errCount
	}
	return int(userCount), errCount
}

func (a *AuthPostgresRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	user, err := a.q.GetUserById(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		t := user.UpdatedAt.Time.UTC()
		updatedAt = &t
	} else {
		updatedAt = nil
	}
	return domain.User{
		ID:           user.ID.Bytes,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.Password,
		IsVerified:   user.IsVerified.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    updatedAt,
	}, nil
}

func (a *AuthPostgresRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	err := a.q.VerifyUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	return err
}

func (a *AuthPostgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := a.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	var updatedAt *time.Time
	if user.UpdatedAt.Valid {
		updatedAt = &user.UpdatedAt.Time
	}
	return domain.User{
		ID:           user.ID.Bytes,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.Password,
		IsVerified:   user.IsVerified.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    updatedAt,
	}, nil
}

func (a *AuthPostgresRepository) UpdatePassword(ctx context.Context, passwordHash string, userID uuid.UUID) error {
	err := a.q.UpdatePassword(ctx, generated.UpdatePasswordParams{
		Password: passwordHash,
		ID:       pgtype.UUID{Bytes: userID, Valid: true},
	})
	return err
}
