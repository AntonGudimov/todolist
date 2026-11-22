package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"todolist/internal/domain"

	"github.com/stretchr/testify/mock"
)

func (h *AuthHandlerSuite) TestVerifyOk() {
	token := "test token"
	h.mockAuthService.On(
		"Verify",
		mock.Anything,
		token,
	).Return(
		nil,
	).Once()

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/v1/auth/verify?token=%v", url.QueryEscape(token)),
		nil,
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusOK, rec.Code)
	h.mockAuthService.AssertExpectations(h.T())
}

func (h *AuthHandlerSuite) TestVerifyUserErr() {
	tests := []struct {
		name               string
		authServiceErr     error
		expectedStatusCode int
	}{
		{
			"email verification not found",
			domain.ErrEmailVerificationNotFound,
			http.StatusNotFound,
		},
		{
			"email verification expired",
			domain.ErrEmailVerificationExpired,
			http.StatusBadRequest,
		},
		{
			"email verification used",
			domain.ErrEmailVerificationUsed,
			http.StatusBadRequest,
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
				token := "test token"
				h.mockAuthService.On(
					"Verify",
					mock.Anything,
					token,
				).Return(
					tt.authServiceErr,
				).Once()

				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/v1/auth/verify?token=%v", url.QueryEscape(token)),
					nil,
				)
				req.Header.Set("Content-Type", "usecasejson")
				rec := httptest.NewRecorder()
				h.router.ServeHTTP(rec, req)

				h.Assert().Equal(tt.expectedStatusCode, rec.Code)
				h.mockAuthService.AssertExpectations(h.T())
			},
		)
	}
}
