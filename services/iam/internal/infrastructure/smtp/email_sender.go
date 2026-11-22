package smtp

import (
	"context"
	"fmt"
	"net/smtp"
	"todolist/internal/configs"
	"todolist/internal/domain"
)

type SMTPSender struct {
	config configs.EmailConfig
	auth   smtp.Auth
}

func NewSMTPSender(config configs.EmailConfig) domain.EmailSender {
	return &SMTPSender{
		config: config,
		auth: smtp.PlainAuth(
			"",
			config.EmailUsername,
			config.EmailPassword,
			config.EmailHost,
		),
	}
}

func (s *SMTPSender) SendVerificationEmail(
	ctx context.Context,
	to string,
	verifyUrl string,
) error {
	subject := "Verify your email"
	body := fmt.Sprintf(
		VerifyEmail,
		verifyUrl,
	)
	message := []byte(
		"From: " + s.config.EmailFrom + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(
		s.config.GetAddr(),
		s.auth,
		s.config.EmailFrom,
		[]string{to},
		message,
	)
}

func (s *SMTPSender) SendForgotPasswordEmail(
	ctx context.Context,
	to string,
	tokenTTL int,
	passwordResetUrl string,
) error {
	subject := "Reset your password"
	body := fmt.Sprintf(
		ForgotPasswordEmail,
		passwordResetUrl,
		tokenTTL,
	)
	message := []byte(
		"From: " + s.config.EmailFrom + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(
		s.config.GetAddr(),
		s.auth,
		s.config.EmailFrom,
		[]string{to},
		message,
	)
}
