package jobs

import (
	"context"
	"log/slog"
	"todolist/internal/configs"
	"todolist/internal/scheduler"
	"todolist/internal/usecase"
)

func GetEmailVerificationCleanUpJob(
	ctx context.Context,
	emailVerificationCleanup *usecase.EmailVerificationCleanup,
	config configs.EmailVerifcationCleanUpConfig,
	logger *slog.Logger,
) scheduler.Job {
	jobName := "email verification cleaner"
	return scheduler.Job{
		Name:     jobName,
		Logger:   logger.With(slog.String("job", jobName)),
		Interval: config.Interval,
		Run:      emailVerificationCleanup.Execute,
	}
}
