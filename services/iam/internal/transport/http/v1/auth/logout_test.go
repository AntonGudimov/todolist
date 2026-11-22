package auth

import (
	"net/http"
	"net/http/httptest"
	"todolist/internal/usecase"

	"github.com/stretchr/testify/mock"
)

func (h *AuthHandlerSuite) TestLogoutOk() {
	refreshToken := "test-refresh-token"
	h.mockAuthService.On(
		"Logout",
		mock.Anything,
		usecase.LogoutCommand{RefreshToken: refreshToken},
	).Return(
		nil,
	).Once()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/logout",
		nil,
	)
	req.AddCookie(
		&http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
		},
	)
	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusNoContent, rec.Code)
	h.Assert().NotContains(rec.Result().Cookies(), "refresh_token")
}
