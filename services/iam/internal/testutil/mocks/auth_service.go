package mocks

import (
	"context"
	"todolist/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) Register(
	ctx context.Context,
	cmd usecase.CreateUserCommand,
) (usecase.CreateUserResult, error) {
	args := m.Called(ctx, cmd)

	result, _ := args.Get(0).(usecase.CreateUserResult)
	return result, args.Error(1)
}

func (m *MockAuthService) Verify(
	ctx context.Context,
	token string,
) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) Login(
	ctx context.Context,
	cmd usecase.LoginUserCommand,
) (usecase.LoginUserResult, error) {
	args := m.Called(ctx, cmd)
	result, _ := args.Get(0).(usecase.LoginUserResult)
	return result, args.Error(1)
}

func (m *MockAuthService) Refresh(
	ctx context.Context,
	cmd usecase.RefreshTokenCommand,
) (usecase.RefreshTokenResult, error) {
	args := m.Called(ctx, cmd)
	result, _ := args.Get(0).(usecase.RefreshTokenResult)
	return result, args.Error(1)
}

func (m *MockAuthService) Logout(
	ctx context.Context,
	cmd usecase.LogoutCommand,
) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(
	ctx context.Context,
	cmd usecase.ForgotPasswordCommand,
) (usecase.ForgotPasswordResult, error) {
	args := m.Called(ctx, cmd)
	result, _ := args.Get(0).(usecase.ForgotPasswordResult)
	return result, args.Error(1)
}

func (m *MockAuthService) ResetPassword(
	ctx context.Context,
	cmd usecase.ResetPasswordCommand,
) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}
