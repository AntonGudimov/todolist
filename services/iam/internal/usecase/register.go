package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserCommand struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResult struct {
	UserID              uuid.UUID
	EmailVerificationID uuid.UUID
	Token               string
}

func (s *Service) Register(
	ctx context.Context,
	cmd CreateUserCommand,
) (CreateUserResult, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(cmd.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return CreateUserResult{}, err
	}
	var result CreateUserResult
	do_err := s.uow.Do(
		ctx,
		func(uow domain.UnitOfWork) error {
			registerResult, registerTxErr := s.registerTx(
				ctx,
				cmd,
				uow,
				hash,
			)
			if registerTxErr != nil {
				return registerTxErr
			}
			result = registerResult
			return nil
		},
	)
	if do_err != nil {
		return CreateUserResult{}, do_err
	}
	sendEmailErr := s.emailSender.SendVerificationEmail(
		ctx,
		cmd.Email,
		fmt.Sprintf(
			"%s/verify?token=%s",
			s.config.AuthPublicUrl,
			url.QueryEscape(result.Token),
		),
	)
	if sendEmailErr != nil {
		return CreateUserResult{}, sendEmailErr
	}
	return result, nil
}

func (s *Service) registerTx(
	ctx context.Context,
	cmd CreateUserCommand,
	uow domain.UnitOfWork,
	hash []byte,
) (CreateUserResult, error) {
	userCount, errCount := s.uow.Auth().CountUsersByEmail(ctx, cmd.Email)
	if errCount != nil {
		return CreateUserResult{}, errCount
	}
	if userCount > 0 {
		return CreateUserResult{}, fmt.Errorf("create user with email: %s: %w", cmd.Email, domain.ErrUserAlreadyExists)
	}
	userID, err := uow.Auth().Register(
		ctx,
		domain.CreateUserParams{
			ID:           uuid.New(),
			Name:         cmd.Name,
			Email:        cmd.Email,
			PasswordHash: string(hash),
		},
	)
	if err != nil {
		return CreateUserResult{}, err
	}
	token, tokenErr := generateToken()
	if tokenErr != nil {
		return CreateUserResult{}, tokenErr
	}
	hashedToken := hashToken(token)
	emailVerification, emailErr := uow.Auth().CreateEmailVerification(
		ctx,
		domain.CreateEmailVerificationParams{
			ID:        uuid.New(),
			UserID:    userID,
			Token:     hashedToken,
			ExpiresAt: time.Now().UTC().Add(s.config.EmailVerificationTTL),
		},
	)
	if emailErr != nil {
		return CreateUserResult{}, emailErr
	}
	return CreateUserResult{UserID: userID, EmailVerificationID: emailVerification.ID, Token: token}, nil
}
