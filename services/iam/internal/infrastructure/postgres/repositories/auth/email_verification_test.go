package auth

import (
	"math"
	"time"
	"todolist/internal/domain"

	"github.com/google/uuid"
)

func (s *AuthPostgresRepositorySuite) TestCreateEmailVerification() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	createEmailVerificationParams := domain.CreateEmailVerificationParams{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "Test token",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	createEmailVerificationResult, err := s.repo.CreateEmailVerification(
		s.ctx,
		createEmailVerificationParams,
	)
	s.Require().NoError(err)

	s.Assert().Equal(createEmailVerificationParams.ID, createEmailVerificationResult.ID)
	s.Assert().Equal(createEmailVerificationParams.UserID, createEmailVerificationResult.UserID)
	s.Assert().Equal(createEmailVerificationParams.Token, createEmailVerificationResult.Token)
	s.Assert().WithinDuration(
		createEmailVerificationParams.ExpiresAt,
		createEmailVerificationResult.ExpiresAt,
		time.Microsecond,
	)
	s.Assert().Nil(createEmailVerificationResult.UsedAt)
}

func (s *AuthPostgresRepositorySuite) TestGetEmailVerificationByTokenOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	createEmailVerificationParams := domain.CreateEmailVerificationParams{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "Test token",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	_, createEmailVerificationErr := s.repo.CreateEmailVerification(
		s.ctx,
		createEmailVerificationParams,
	)
	s.Require().NoError(createEmailVerificationErr)

	emailVerification, err := s.repo.GetEmailVerificationByToken(
		s.ctx, createEmailVerificationParams.Token,
	)
	s.Require().NoError(err)

	s.Assert().Equal(createEmailVerificationParams.ID, emailVerification.ID)
	s.Assert().Equal(createEmailVerificationParams.UserID, emailVerification.UserID)
	s.Assert().Equal(createEmailVerificationParams.Token, emailVerification.Token)
	s.Assert().WithinDuration(
		createEmailVerificationParams.ExpiresAt,
		emailVerification.ExpiresAt,
		time.Millisecond,
	)
	s.Assert().Nil(emailVerification.UsedAt)
}

func (s *AuthPostgresRepositorySuite) TestGetEmailVerificationByTokenNotFound() {
	emailVerification, err := s.repo.GetEmailVerificationByToken(s.ctx, "test token")
	s.Assert().ErrorIs(err, domain.ErrEmailVerificationNotFound)
	s.Assert().Equal(domain.EmailVerification{}, emailVerification)
}

func (s *AuthPostgresRepositorySuite) TestUseEmailVerification() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	createEmailVerificationParams := domain.CreateEmailVerificationParams{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "Test token",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	_, createEmailVerificationErr := s.repo.CreateEmailVerification(
		s.ctx,
		createEmailVerificationParams,
	)
	s.Require().NoError(createEmailVerificationErr)

	usedAt := time.Now().UTC()
	err := s.repo.UseEmailVerification(s.ctx, domain.UseEmailVerificationParams{ID: createEmailVerificationParams.ID, UsedAt: usedAt})
	s.Require().NoError(err)

	emailVerification, getEmailErr := s.repo.GetEmailVerificationByToken(s.ctx, createEmailVerificationParams.Token)
	s.Require().NoError(getEmailErr)

	s.Assert().WithinDuration(
		usedAt,
		*emailVerification.UsedAt,
		time.Millisecond,
	)
}

func (s *AuthPostgresRepositorySuite) TestCleanupEmailVerificationsByCutoffOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	retention := 2 * time.Hour
	createEmailVerificationParams := domain.CreateEmailVerificationParams{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "Test token",
		ExpiresAt: time.Now().UTC().Add(-retention),
	}
	_, createEmailVerificationErr := s.repo.CreateEmailVerification(
		s.ctx,
		createEmailVerificationParams,
	)
	s.Require().NoError(createEmailVerificationErr)

	_, err := s.repo.CleanupEmailVerifications(s.ctx, createEmailVerificationParams.ExpiresAt, 1)
	s.Require().NoError(err)

	emailVerification, emailVerificationErr := s.repo.GetEmailVerificationByToken(s.ctx, createEmailVerificationParams.Token)
	s.Assert().Equal(domain.EmailVerification{}, emailVerification)
	s.Assert().ErrorIs(emailVerificationErr, domain.ErrEmailVerificationNotFound)
}

func (s *AuthPostgresRepositorySuite) TestCleanupEmailVerificationsByUsedAtOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	retention := 2 * time.Hour
	createEmailVerificationParams := domain.CreateEmailVerificationParams{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "Test token",
		ExpiresAt: time.Now().UTC().Add(retention),
	}
	_, createEmailVerificationErr := s.repo.CreateEmailVerification(
		s.ctx,
		createEmailVerificationParams,
	)
	s.Require().NoError(createEmailVerificationErr)

	useEmailVerificationErr := s.repo.UseEmailVerification(
		s.ctx,
		domain.UseEmailVerificationParams{ID: createEmailVerificationParams.ID, UsedAt: time.Now().UTC()},
	)
	s.Require().NoError(useEmailVerificationErr)

	_, err := s.repo.CleanupEmailVerifications(s.ctx, time.Now().UTC(), 1)
	s.Require().NoError(err)

	emailVerification, emailVerificationErr := s.repo.GetEmailVerificationByToken(s.ctx, createEmailVerificationParams.Token)
	s.Assert().Equal(domain.EmailVerification{}, emailVerification)
	s.Assert().ErrorIs(emailVerificationErr, domain.ErrEmailVerificationNotFound)
}

func (s *AuthPostgresRepositorySuite) TestCleanupEmailVerificationsByBatchOk() {
	createUserParams := domain.CreateUserParams{
		ID:           uuid.New(),
		Name:         "Test name",
		Email:        "Test email",
		PasswordHash: "Test password",
	}
	userID, registerErr := s.repo.Register(
		s.ctx, createUserParams,
	)
	s.Require().NoError(registerErr)

	retention := 2 * time.Hour
	expiresAt := time.Now().UTC().Add(-retention)
	emailVerificationArrayLen := 5
	for range emailVerificationArrayLen {
		createEmailVerificationParams := domain.CreateEmailVerificationParams{
			ID:        uuid.New(),
			UserID:    userID,
			Token:     uuid.NewString(),
			ExpiresAt: expiresAt,
		}
		_, createEmailVerificationErr := s.repo.CreateEmailVerification(
			s.ctx,
			createEmailVerificationParams,
		)
		s.Require().NoError(createEmailVerificationErr)
	}
	batchSize := 2
	expectedDeletedEmailVerificationsCountArr := []int{batchSize, batchSize, batchSize - 1}

	actualDeletedEmailVerificationsCountArr := make([]int, 0, emailVerificationArrayLen)
	for range int(math.Ceil(float64(emailVerificationArrayLen) / float64(2))) {
		ids, err := s.repo.CleanupEmailVerifications(s.ctx, expiresAt, batchSize)
		s.Require().NoError(err)
		actualDeletedEmailVerificationsCountArr = append(actualDeletedEmailVerificationsCountArr, len(ids))
	}
	s.Assert().Equal(expectedDeletedEmailVerificationsCountArr, actualDeletedEmailVerificationsCountArr)
}
