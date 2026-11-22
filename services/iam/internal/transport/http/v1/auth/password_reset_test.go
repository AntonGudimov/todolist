package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/stretchr/testify/mock"
)

func (h *AuthHandlerSuite) TestForgotPasswordOk() {
	forgotRequest := ForgotRequest{Email: "testemail@yandex.ru"}
	token := "test token"
	h.mockAuthService.On(
		"ForgotPassword",
		mock.Anything,
		usecase.ForgotPasswordCommand{
			Email: forgotRequest.Email,
		},
	).Return(
		usecase.ForgotPasswordResult{Token: token},
		nil,
	).Once()

	body, bodyErr := json.Marshal(forgotRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/password/forgot",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")

	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	var response map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	h.Require().NoError(err)

	h.Assert().Equal(http.StatusOK, rec.Code)
	h.Assert().Equal(token, response["token"])
}

func (h *AuthHandlerSuite) TestForgotPasswordInvalidInputErr() {
	forgotRequest := map[string]string{
		"Email": "invalid email",
	}
	body, bodyErr := json.Marshal(forgotRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/password/forgot",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)
	h.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (h *AuthHandlerSuite) TestForgotPasswordErr() {
	tests := []struct {
		name           string
		authServiceErr error
		statusCode     int
	}{
		{
			"user not found",
			domain.ErrUserNotFound,
			http.StatusNotFound,
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
				forgotRequest := ForgotRequest{Email: "testemail@yandex.ru"}
				h.mockAuthService.On(
					"ForgotPassword",
					mock.Anything,
					usecase.ForgotPasswordCommand{
						Email: forgotRequest.Email,
					},
				).Return(
					nil,
					tt.authServiceErr,
				).Once()

				body, bodyErr := json.Marshal(forgotRequest)
				h.Require().NoError(bodyErr)

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/auth/password/forgot",
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "usecasejson")

				rec := httptest.NewRecorder()
				h.router.ServeHTTP(rec, req)

				h.Assert().Equal(tt.statusCode, rec.Code)
			},
		)
	}
}

func (h *AuthHandlerSuite) TestResetPasswordOk() {
	resetPasswordRequest := ResetPasswordRequest{
		Token:       "test token",
		NewPassword: "test password",
	}
	h.mockAuthService.On(
		"ResetPassword",
		mock.Anything,
		usecase.ResetPasswordCommand{
			Token:       resetPasswordRequest.Token,
			NewPassword: resetPasswordRequest.NewPassword,
		},
	).Return(
		nil,
	).Once()

	body, bodyErr := json.Marshal(resetPasswordRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/password/reset",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")

	rec := httptest.NewRecorder()
	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (h *AuthHandlerSuite) TestResetPasswordInvalidInputErr() {
	resetPasswordRequest := map[string]string{}
	body, bodyErr := json.Marshal(resetPasswordRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/password/reset",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)
	h.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (h *AuthHandlerSuite) TestResetPasswordErr() {
	tests := []struct {
		name           string
		authServiceErr error
		statusCode     int
	}{
		{
			"password reset token not found",
			domain.ErrPasswordResetTokenNotFound,
			http.StatusNotFound,
		},
		{
			"password reset token revoked",
			domain.ErrPasswordResetTokenRevoked,
			http.StatusBadRequest,
		},
		{
			"password reset token used",
			domain.ErrPasswordResetTokenUsed,
			http.StatusBadRequest,
		},
		{
			"password reset token expired",
			domain.ErrPasswordResetTokenExpired,
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
				resetPasswordRequest := ResetPasswordRequest{
					Token:       "test token",
					NewPassword: "test password",
				}
				h.mockAuthService.On(
					"ResetPassword",
					mock.Anything,
					usecase.ResetPasswordCommand{
						Token:       resetPasswordRequest.Token,
						NewPassword: resetPasswordRequest.NewPassword,
					},
				).Return(
					tt.authServiceErr,
				).Once()

				body, bodyErr := json.Marshal(resetPasswordRequest)
				h.Require().NoError(bodyErr)

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/auth/password/reset",
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "usecasejson")

				rec := httptest.NewRecorder()
				h.router.ServeHTTP(rec, req)

				h.Assert().Equal(tt.statusCode, rec.Code)
			},
		)
	}
}
