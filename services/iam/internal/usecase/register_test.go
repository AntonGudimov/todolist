package usecase

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthServiceSuite) TestRegisterOk() {
	createUserCommand := CreateUserCommand{
		Name:     "test name",
		Email:    "test email",
		Password: "test password",
	}

	registerResult, registerErr := s.service.Register(
		s.ctx,
		createUserCommand,
	)
	s.Require().NoError(registerErr)

	s.Assert().Len(s.uow.Db.Users, 1)
	s.Assert().Len(s.uow.Db.EmailVerifications, 1)
	s.Assert().Len(s.emailSender.SentVerificationEmails, 1)

	user := s.uow.Db.Users[registerResult.UserID]

	emailVerification := s.uow.Db.EmailVerifications[registerResult.EmailVerificationID]

	// assert user
	hashErr := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(createUserCommand.Password),
	)
	s.Require().NoError(hashErr)
	s.Assert().Equal(createUserCommand.Name, user.Name)
	s.Assert().Equal(createUserCommand.Email, user.Email)
	s.Assert().Equal(false, user.IsVerified)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		user.CreatedAt,
		time.Second,
	)
	s.Assert().Nil(user.UpdatedAt)

	// assert email verification
	s.Assert().Equal(emailVerification.UserID, user.ID)
	s.Assert().NotEmpty(emailVerification.Token)
	s.Assert().WithinDuration(
		emailVerification.ExpiresAt,
		time.Now().UTC().Add(s.config.EmailVerificationTTL),
		time.Second,
	)
	s.Assert().Nil(emailVerification.UsedAt)
}

func (s *AuthServiceSuite) TestRegisterUserAlreadyExistErr() {
	createUserCommand := CreateUserCommand{
		Name:     "test name",
		Email:    "test email",
		Password: "test password",
	}
	s.service.Register(
		s.ctx, createUserCommand,
	)

	registerResult, err := s.service.Register(
		s.ctx, createUserCommand,
	)

	s.Require().ErrorIs(err, domain.ErrUserAlreadyExists)
	s.Assert().Equal(uuid.Nil, registerResult.UserID)
}
