package jobs

import (
	"context"
	"log/slog"
	"todolist/internal/configs"
	"todolist/internal/domain"
	"todolist/internal/scheduler"
)

func GetKeyReloaderJob(ctx context.Context, keystore domain.KeyStore, config configs.KeyReloaderConfig, logger *slog.Logger) scheduler.Job {
	jobName := "key reloader"
	return scheduler.Job{
		Name:     jobName,
		Logger:   logger.With(slog.String("job", jobName)),
		Interval: config.Interval,
		Run: func(ctx context.Context) error {
			return keystore.Reload()
		},
	}
}
