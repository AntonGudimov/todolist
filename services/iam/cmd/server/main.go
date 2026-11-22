package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"todolist/internal/application"
	"todolist/internal/closer"
	"todolist/internal/configs"
	"todolist/internal/infrastructure/jwt"
	"todolist/internal/infrastructure/postgres"
	"todolist/internal/infrastructure/smtp"
	"todolist/internal/logging"
	"todolist/internal/scheduler"
	"todolist/internal/scheduler/jobs"
	"todolist/internal/transport/http"
	"todolist/internal/transport/http/health"
	v1Auth "todolist/internal/transport/http/v1/auth"
	"todolist/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, configErr := configs.NewConfig()
	if configErr != nil {
		log.Fatalf("Couldn't initialize config: %v", configErr)
	}
	logger, loggerErr := logging.GetLogger(config.Logger)
	if loggerErr != nil {
		log.Fatalf("Couldn't initialize logger: %v", loggerErr)
	}
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()
	cleanup := &closer.Cleanup{}

	dbPool, err := pgxpool.New(ctx, config.Database.GetURL())
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	cleanup.Add(func() error {
		logger.Info("Closing db pool", "component", "cleanup")
		dbPool.Close()
		return nil
	})

	uow := postgres.NewPostgres(dbPool)
	keyStore, keyStoreErr := jwt.NewFileKeyStore(config.KeyStore.Dir)
	if keyStoreErr != nil {
		log.Fatalf("Error creating key store: %v", keyStoreErr)
	}
	jwtService := jwt.NewJWTService(keyStore, &config.JWT)
	smtpSender := smtp.NewSMTPSender(config.Email)
	authService := usecase.NewAuthService(
		uow,
		smtpSender,
		jwtService,
		config.Auth,
	)
	authHandlerV1 := v1Auth.NewAuthHandler(authService, logger)
	healthHandler := health.NewHealthHandler()

	emailVerificationCleanUp := usecase.NewEmailVerificationCleanup(uow, config.EmailVerifcationCleanup, logger)

	httpApp := http.NewApp(config.Host, logger, authHandlerV1, healthHandler, config.ShutdownTimeout/3)
	jobs := []scheduler.Job{
		jobs.GetKeyReloaderJob(ctx, keyStore, config.KeyReloader, logger),
		jobs.GetEmailVerificationCleanUpJob(ctx, emailVerificationCleanUp, config.EmailVerifcationCleanup, logger),
	}
	scheduler := scheduler.NewScheduler(jobs, logger)
	components := []application.Component{
		httpApp, scheduler,
	}
	app := application.New(components, cleanup, logger, config.ShutdownTimeout)
	runErr := app.Run(ctx)
	if runErr != nil {
		logger.Error(runErr.Error())
	}
}
