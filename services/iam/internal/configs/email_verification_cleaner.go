package configs

import "time"

type EmailVerifcationCleanUpConfig struct {
	Interval  time.Duration `env:"EMAIL_VERIFICATION_CLEAN_UP_INTERVAL"`
	BatchSize int           `env:"EMAIL_VERIFICATION_CLEAN_UP_BATCH_SIZE"`
}
