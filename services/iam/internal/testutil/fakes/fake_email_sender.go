package fakes

import (
	"context"
)

type EmailVerificationInfo struct {
	VerifyUrl string
}

type PasswordResetInfo struct {
	ResetUrl string
}

type FakeEmailSender struct {
	SentVerificationEmails  []EmailVerificationInfo
	SentPasswordResetTokens []PasswordResetInfo
}

func NewFakeEmailSender() *FakeEmailSender {
	return &FakeEmailSender{
		SentVerificationEmails:  make([]EmailVerificationInfo, 0),
		SentPasswordResetTokens: make([]PasswordResetInfo, 0),
	}
}

func (s *FakeEmailSender) SendVerificationEmail(ctx context.Context, to string, verifyUrl string) error {
	s.SentVerificationEmails = append(s.SentVerificationEmails, EmailVerificationInfo{VerifyUrl: verifyUrl})
	return nil
}

func (s *FakeEmailSender) SendForgotPasswordEmail(ctx context.Context, to string, tokenTTL int, passwordResetUrl string) error {
	s.SentPasswordResetTokens = append(s.SentPasswordResetTokens, PasswordResetInfo{ResetUrl: passwordResetUrl})
	return nil
}
