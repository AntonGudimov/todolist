package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"todolist/internal/configs"
	"todolist/internal/transport/http/health"
	"todolist/internal/transport/http/middleware"
	v1 "todolist/internal/transport/http/v1"
	v1Auth "todolist/internal/transport/http/v1/auth"

	"github.com/gin-gonic/gin"
)

type App struct {
	logger          *slog.Logger
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewApp(
	config configs.HostConfig,
	logger *slog.Logger,
	authHandlerV1 *v1Auth.AuthHandler,
	healthHandler *health.HealthHandler,
	shutdownTimeout time.Duration,
) *App {
	router := gin.New()
	registerRouters(router, logger, authHandlerV1, healthHandler)
	return &App{
		logger: logger.With(slog.String("component", "http server")),
		server: &http.Server{
			Addr:    config.GetUrl(),
			Handler: router.Handler(),
		},
		shutdownTimeout: shutdownTimeout,
	}
}

func (app *App) Run(ctx context.Context) error {
	app.logger.Info("starting up")
	err := app.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}

func (app *App) Shutdown(ctx context.Context) error {
	app.logger.Info("shutting down")
	err := app.server.Shutdown(ctx)
	return err
}

func registerRouters(
	router *gin.Engine,
	logger *slog.Logger,
	authHandlerV1 *v1Auth.AuthHandler,
	healthHandler *health.HealthHandler,
) {
	router.Use(gin.Recovery())

	router.GET("/health", healthHandler.Health)

	api := router.Group("/api")
	api.Use(middleware.RequestIDMiddleware())
	api.Use(middleware.SlogMiddleware(logger))

	v1RouterGroup := api.Group("/v1")
	v1.RegisterRouters(v1RouterGroup, authHandlerV1)
}
