package auth

import (
	"context"
	"testing"
	"todolist/internal/configs"
	"todolist/internal/logging"
	"todolist/internal/testutil/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerSuite struct {
	suite.Suite

	ctx             context.Context
	mockAuthService *mocks.MockAuthService
	handler         AuthHandler
	router          *gin.Engine
}

func (h *AuthHandlerSuite) SetupTest() {
	mockAuthService := mocks.NewMockAuthService()
	logger, _ := logging.GetLogger(configs.LoggerConfig{Level: "info"})
	authHandler := NewAuthHandler(mockAuthService, logger)
	router := gin.New()
	router.POST(
		"/api/v1/auth/register",
		authHandler.CreateUser,
	)
	router.GET(
		"/api/v1/auth/verify",
		authHandler.VerifyUser,
	)
	router.POST(
		"/api/v1/auth/login",
		authHandler.Login,
	)
	router.POST(
		"/api/v1/auth/refresh",
		authHandler.Refresh,
	)
	router.POST(
		"/api/v1/auth/logout",
		authHandler.Logout,
	)
	router.POST(
		"/api/v1/auth/password/forgot",
		authHandler.ForgotPassword,
	)
	router.POST(
		"/api/v1/auth/password/reset",
		authHandler.ResetPassword,
	)
	h.ctx = context.Background()
	h.mockAuthService = mockAuthService
	h.handler = *authHandler
	h.router = router
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
