package configs

import "time"

type KeyReloaderConfig struct {
	Interval time.Duration `env:"KEY_RELOADER_INTERVAL"`
}
