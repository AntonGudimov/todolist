package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"todolist/internal/configs"
	"todolist/internal/domain"
)

type EmailVerificationCleanup struct {
	uow    domain.UnitOfWork
	config configs.EmailVerifcationCleanUpConfig
	logger *slog.Logger
}

func NewEmailVerificationCleanup(
	uow domain.UnitOfWork,
	config configs.EmailVerifcationCleanUpConfig,
	logger *slog.Logger,
) *EmailVerificationCleanup {
	return &EmailVerificationCleanup{
		uow:    uow,
		config: config,
		logger: logger,
	}
}

func (evc *EmailVerificationCleanup) Execute(ctx context.Context) error {
	utcNow := time.Now().UTC()
	idsStr := make([]string, 0)
	for {
		ids, err := evc.uow.Auth().CleanupEmailVerifications(ctx, utcNow, evc.config.BatchSize)
		if err != nil {
			return fmt.Errorf("cleanup email verifications: %w", err)
		}
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			idsStr = append(idsStr, id.String())
		}

	}
	evc.logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Cleaned up email verifications",
		slog.Int("count", len(idsStr)),
		slog.String("ids", strings.Join(idsStr, ", ")),
	)
	return nil
}
