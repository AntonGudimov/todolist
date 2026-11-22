package configs

import (
	"time"
)

type JWTConfig struct {
	Issuer   string        `env:"JWT_ISSUER"`
	TokenTTL time.Duration `env:"JWT_TOKEN_TTL"`
}
