package usecase

import (
	"context"
	"testing"
	"time"
	"todolist/internal/configs"
	"todolist/internal/testutil/fakes"

	"github.com/stretchr/testify/suite"
)

type AuthServiceSuite struct {
	suite.Suite

	uow         *fakes.FakeUnitOfWork
	config      configs.AuthConfig
	service     AuthService
	emailSender *fakes.FakeEmailSender
	jwtService  *fakes.FakeJWTService

	ctx context.Context
}

func (s *AuthServiceSuite) SetupTest() {
	s.ctx = context.Background()
	db := fakes.NewFakeDB()
	uow := fakes.NewFakeUnitOfWork(db)
	config := configs.AuthConfig{
		EmailVerificationTTL: time.Hour,
		AuthPublicUrl:        "http://iam.lan/auth",
		RefreshTokenTTL:      720 * time.Hour,
	}
	emailSender := fakes.NewFakeEmailSender()
	jwtService := fakes.NewFakeJWTService("test access token", nil)
	s.uow = uow
	s.config = config
	s.emailSender = emailSender
	s.jwtService = jwtService
	s.service = NewAuthService(uow, emailSender, jwtService, config)
}

func (s *AuthServiceSuite) TearDownTest() {
	s.uow.Db.Clear()
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
