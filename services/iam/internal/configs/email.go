package configs

import (
	"fmt"
	"time"
)

type EmailConfig struct {
	EmailHost     string `env:"EMAIL_HOST"`
	EmailPort     int    `env:"EMAIL_PORT"`
	EmailUsername string `env:"EMAIL_USERNAME"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
	EmailFrom     string `env:"EMAIL_FROM"`

	EmailVerificationTTL time.Duration `env:"EMAIL_VERIFICATION_TTL"`
}

func (c *EmailConfig) GetAddr() string {
	return fmt.Sprintf(
		"%v:%v",
		c.EmailHost,
		c.EmailPort,
	)
}
