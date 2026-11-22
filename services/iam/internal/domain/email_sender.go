package domain

import (
	"context"
)

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, to string, verifyURL string) error
	SendForgotPasswordEmail(ctx context.Context, to string, tokenTTL int, passwordResetURL string) error
}
