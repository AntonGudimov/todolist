package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	healthHandler := NewHealthHandler()
	router := gin.New()
	router.GET("/api/health", healthHandler.Health)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/health",
		nil,
	)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
