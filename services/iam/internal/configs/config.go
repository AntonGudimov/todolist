package configs

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ShutdownTimeout         time.Duration `env:"SHUTDOWN_TIMEOUT"`
	Database                DatabaseConfig
	Host                    HostConfig
	Email                   EmailConfig
	Auth                    AuthConfig
	JWT                     JWTConfig
	KeyStore                KeyStoreConfig
	KeyReloader             KeyReloaderConfig
	EmailVerifcationCleanup EmailVerifcationCleanUpConfig
	Logger                  LoggerConfig
}

func NewConfig() (*Config, error) {
	var config Config

	err := cleanenv.ReadEnv(&config)
	return &config, err
}
