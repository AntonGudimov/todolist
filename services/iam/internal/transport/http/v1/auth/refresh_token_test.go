package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/stretchr/testify/mock"
)

func (h *AuthHandlerSuite) TestRefreshOk() {
	oldRefreshToken := "old-refresh-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"
	newRefreshTokenExpiresAt := time.Now().Add(time.Hour * 720)

	h.mockAuthService.On(
		"Refresh",
		mock.Anything,
		usecase.RefreshTokenCommand{
			RefreshToken: oldRefreshToken,
		},
	).Return(
		usecase.RefreshTokenResult{
			AccessToken:           newAccessToken,
			RefreshToken:          newRefreshToken,
			RefreshTokenExpiresAt: newRefreshTokenExpiresAt,
		}, nil,
	).Once()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/refresh",
		nil,
	)
	req.AddCookie(
		&http.Cookie{
			Name:  "refresh_token",
			Value: oldRefreshToken,
		},
	)
	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	var response map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	h.Require().NoError(err)

	var refreshTokenCookie *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "refresh_token" {
			refreshTokenCookie = cookie
			break
		}
	}

	h.Assert().Equal(http.StatusOK, rec.Code)
	h.Assert().Equal("Bearer", response["type"])
	h.Assert().Equal(newAccessToken, response["accessToken"])

	h.Assert().Equal(newRefreshToken, refreshTokenCookie.Value)
	h.Assert().Equal(int(time.Until(newRefreshTokenExpiresAt).Seconds()), refreshTokenCookie.MaxAge)
	h.Assert().Equal("/auth", refreshTokenCookie.Path)
	h.Assert().Equal("", refreshTokenCookie.Domain)
	h.Assert().Equal(true, refreshTokenCookie.Secure)
	h.Assert().Equal(true, refreshTokenCookie.HttpOnly)

	h.mockAuthService.AssertExpectations(h.T())

}

func (h *AuthHandlerSuite) TestRefreshNoCookieErr() {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/refresh",
		nil,
	)
	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusUnauthorized, rec.Code)

	h.mockAuthService.AssertExpectations(h.T())
}

func (h *AuthHandlerSuite) TestRefreshErr() {
	tests := []struct {
		name           string
		authServiceErr error
		statusCode     int
	}{
		{
			"refresh token not found",
			domain.ErrRefreshTokenNotFound,
			http.StatusUnauthorized,
		},
		{
			"refresh token revoked",
			domain.ErrRefreshTokenInvalid,
			http.StatusUnauthorized,
		},
		{
			"refresh token expired",
			domain.ErrRefreshTokenInvalid,
			http.StatusUnauthorized,
		},
		{
			"internal server error",
			errors.New("internal server error"),
			http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		h.Run(
			tt.name,
			func() {
				oldRefreshToken := "old-refresh-token"
				h.mockAuthService.On(
					"Refresh",
					mock.Anything,
					usecase.RefreshTokenCommand{
						RefreshToken: oldRefreshToken,
					},
				).Return(
					usecase.RefreshTokenResult{}, tt.authServiceErr,
				).Once()

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/auth/refresh",
					nil,
				)
				req.AddCookie(
					&http.Cookie{
						Name:  "refresh_token",
						Value: oldRefreshToken,
					},
				)
				rec := httptest.NewRecorder()
				h.router.ServeHTTP(rec, req)

				h.Assert().Equal(tt.statusCode, rec.Code)

				h.mockAuthService.AssertExpectations(h.T())
			},
		)
	}
}
