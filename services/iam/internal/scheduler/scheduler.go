package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Job struct {
	Name     string
	Logger   *slog.Logger
	Interval time.Duration
	Run      func(context.Context) error
}

type Scheduler struct {
	jobs             []Job
	logger           *slog.Logger
	stopSchedulingCh chan struct{}
	done             chan struct{}
	executionCtx     context.Context
	executionCancel  context.CancelFunc
}

func NewScheduler(jobs []Job, logger *slog.Logger) *Scheduler {
	executionCtx, exectionCancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:             jobs,
		logger:           logger.With(slog.String("component", "scheduler")),
		stopSchedulingCh: make(chan struct{}),
		done:             make(chan struct{}),
		executionCtx:     executionCtx,
		executionCancel:  exectionCancel,
	}
}

func (s *Scheduler) Run(parentCtx context.Context) error {
	s.logger.Info("starting up")

	var wg sync.WaitGroup
	for _, job := range s.jobs {
		wg.Go(func() {
			s.runJob(job)
		})
	}
	go func() {
		wg.Wait()
		close(s.done)
	}()
	<-parentCtx.Done()
	return nil
}

func (s *Scheduler) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down")
	close(s.stopSchedulingCh)
	select {
	case <-s.done:
		s.logger.Info("successfully shut down jobs")
		return nil
	case <-ctx.Done():
		s.logger.Error("shutting down jobs execeeded context deadline")
		s.executionCancel()
		return nil
	}
}

func (s *Scheduler) runJob(job Job) {
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopSchedulingCh:
			return

		case <-ticker.C:
			if err := job.Run(s.executionCtx); err != nil {
				s.logger.Error(
					"job failed",
					"job", job.Name,
					"error", err,
				)
			} else {
				s.logger.Info(
					"job succeded",
					"job", job.Name,
				)
			}
		}
	}
}
