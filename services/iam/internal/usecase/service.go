package usecase

import (
	"context"
	"todolist/internal/configs"
	"todolist/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, cmd CreateUserCommand) (CreateUserResult, error)
	Verify(ctx context.Context, token string) error
	Login(ctx context.Context, cmd LoginUserCommand) (LoginUserResult, error)
	Refresh(ctx context.Context, cmd RefreshTokenCommand) (RefreshTokenResult, error)
	Logout(ctx context.Context, cmd LogoutCommand) error
	ForgotPassword(ctx context.Context, cmd ForgotPasswordCommand) (ForgotPasswordResult, error)
	ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error
}

type Service struct {
	uow         domain.UnitOfWork
	emailSender domain.EmailSender
	jwtService  domain.JWTService
	config      configs.AuthConfig
}

func NewAuthService(uow domain.UnitOfWork, emailSender domain.EmailSender, jwtService domain.JWTService, config configs.AuthConfig) AuthService {
	return &Service{uow: uow, emailSender: emailSender, jwtService: jwtService, config: config}
}
