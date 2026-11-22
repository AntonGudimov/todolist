package middleware

import (
	"context"
	"log/slog"
	"time"
	"todolist/internal/constants"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(constants.RequestIDKey, requestID)
		ctx := context.WithValue(
			c.Request.Context(),
			constants.RequestIDKey,
			requestID,
		)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		requestID, _ := c.Get(constants.RequestIDKey)
		logger.Info("http request completed",
			slog.String(string(constants.RequestIDKey), requestID.(string)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}
