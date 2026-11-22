package auth

import (
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthPostgresRepositorySuite) TestRegister() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	actualUserID, registerErr := s.repo.Register(
		s.ctx,
		createUserParams,
	)
	s.Require().NoError(registerErr)

	user, err := s.repo.GetUserByID(s.ctx, actualUserID)
	s.Require().NoError(err)

	s.Assert().Equal(createUserParams.ID, user.ID)
	s.Assert().Equal(createUserParams.Name, user.Name)
	s.Assert().Equal(createUserParams.Email, user.Email)
	s.Assert().Equal(createUserParams.PasswordHash, user.PasswordHash)
	s.Assert().False(user.IsVerified)
}

func (s *AuthPostgresRepositorySuite) TestCountUsersByEmail() {
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

	expectedUserCount := 1
	actualUserCount, err := s.repo.CountUsersByEmail(s.ctx, createUserParams.Email)
	s.Require().NoError(err)
	s.Assert().Equal(expectedUserCount, actualUserCount)
}

func (s *AuthPostgresRepositorySuite) TestUserNotFoundByID() {
	_, err := s.repo.GetUserByID(s.ctx, uuid.New())
	s.Require().ErrorIs(err, domain.ErrUserNotFound)
}

func (s *AuthPostgresRepositorySuite) TestVerifyUser() {
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

	verifyErr := s.repo.VerifyUser(s.ctx, createUserParams.ID)
	s.Require().NoError(verifyErr)

	user, err := s.repo.GetUserByID(s.ctx, createUserParams.ID)
	s.Require().NoError(err)

	s.Assert().Equal(createUserParams.ID, user.ID)
	s.Assert().Equal(true, user.IsVerified)
	s.Assert().WithinDuration(
		time.Now().UTC(),
		*user.UpdatedAt,
		time.Millisecond*5,
	)
}

func (s *AuthPostgresRepositorySuite) TestGetUserByEmailOk() {
	userEmail := "Test email"
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        userEmail,
		PasswordHash: "Test password",
	}
	_, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	user, userErr := s.repo.GetUserByEmail(s.ctx, userEmail)
	s.Require().NoError(userErr)

	s.Assert().Equal(createUserParams.Name, user.Name)
	s.Assert().Equal(createUserParams.Email, user.Email)
	s.Assert().Equal(createUserParams.PasswordHash, user.PasswordHash)
}

func (s *AuthPostgresRepositorySuite) TestGetUserByEmailNotFound() {
	user, userErr := s.repo.GetUserByEmail(s.ctx, "non existend email")
	s.Assert().Equal(domain.User{}, user)
	s.Assert().ErrorIs(userErr, domain.ErrUserNotFound)
}

func (s *AuthPostgresRepositorySuite) TestUpdatePasswordOk() {
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

	newPassword := "new test password"

	err := s.repo.UpdatePassword(s.ctx, newPassword, createUserParams.ID)
	s.Require().NoError(err)

	user, userErr := s.repo.GetUserByID(s.ctx, createUserParams.ID)
	s.Require().NoError(userErr)

	s.Assert().Equal(newPassword, user.PasswordHash)
}
