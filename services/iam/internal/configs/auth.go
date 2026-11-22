package configs

import "time"

type AuthConfig struct {
	EmailVerificationTTL  time.Duration `env:"AUTH_EMAIL_VERIFICATION_TTL"`
	AuthPublicUrl         string        `env:"AUTH_PUBLIC_URL"`
	RefreshTokenTTL       time.Duration `env:"REFRESH_TOKEN_TTL"`
	PasswordResetTokenTTL time.Duration `env:"PASSWORD_RESET_TOKEN_TTL"`
}
