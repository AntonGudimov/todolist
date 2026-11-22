package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"todolist/internal/closer"
)

type App struct {
	logger          *slog.Logger
	components      []Component
	closer          *closer.Cleanup
	shutdownTimeout time.Duration
}

func New(components []Component, closer *closer.Cleanup, logger *slog.Logger, ShutdownTimeout time.Duration) *App {
	return &App{components: components, closer: closer, logger: logger, shutdownTimeout: ShutdownTimeout}
}

func (a *App) Run(parentCtx context.Context) error {
	a.logger.Info("starting up application")
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()
	errCh := make(chan error, len(a.components))

	// start each component
	var wg sync.WaitGroup
	for _, component := range a.components {
		wg.Go(func() {
			if err := component.Run(runCtx); err != nil {
				errCh <- err
			}
		})
	}
	var runtimeErr error
	var shutdownErr error
	var cleanupErr error

	select {
	case <-parentCtx.Done():
		a.logger.Info("received shutdown request")
	case err := <-errCh:
		runtimeErr = errors.Join(err)
	}

	// stop all components
	a.logger.Info("shutting down application")
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(), a.shutdownTimeout,
	)
	defer shutdownCancel()
	shutdownErr = a.shutdownComponents(shutdownCtx)
	runDone := make(chan struct{})
	go func() {
		runCancel()
		wg.Wait()
		close(errCh)
		close(runDone)
	}()

	select {
	case <-runDone:
		// Collect errors from all components.
		for err := range errCh {
			runtimeErr = errors.Join(runtimeErr, err)
		}

		// Close dependencies only after all components have stopped.
		if err := a.closer.Close(shutdownCtx); err != nil {
			cleanupErr = errors.Join(
				cleanupErr,
				fmt.Errorf("close dependencies: %w", err),
			)
		}
		return errors.Join(
			runtimeErr,
			shutdownErr,
			cleanupErr,
		)
	case <-shutdownCtx.Done():
		return errors.Join(runtimeErr, shutdownErr)
	}
}

func (a *App) shutdownComponents(ctx context.Context) error {
	var wg sync.WaitGroup
	var shutdownErr error
	var mu sync.Mutex

	// shutdown the components
	for _, component := range a.components {
		wg.Go(func() {
			if err := component.Shutdown(ctx); err != nil {
				mu.Lock()
				shutdownErr = errors.Join(shutdownErr, err)
				mu.Unlock()
			}
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return shutdownErr
	case <-ctx.Done():
		return errors.Join(
			shutdownErr,
			fmt.Errorf("shutting down components exceeded context deadline: %v", a.shutdownTimeout),
		)
	}
}
