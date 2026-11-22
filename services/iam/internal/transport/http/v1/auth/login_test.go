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

func (h *AuthHandlerSuite) TestLoginUserOk() {
	loginRequest := LoginRequest{
		Email:    "testEmail@google.com",
		Password: "test password",
	}
	accessToken := "test access token"
	refreshToken := "test refresh token"
	h.mockAuthService.On(
		"Login",
		mock.Anything,
		usecase.LoginUserCommand{
			Email:    loginRequest.Email,
			Password: loginRequest.Password,
		},
	).Return(
		usecase.LoginUserResult{AccessToken: accessToken, RefreshToken: refreshToken},
		nil,
	).Once()

	body, bodyErr := json.Marshal(loginRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)

	var response map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	h.Require().NoError(err)

	h.Assert().Equal(http.StatusOK, rec.Code)
	h.Assert().Equal("Bearer", response["type"])
	h.Assert().Equal(accessToken, response["accessToken"])
	h.mockAuthService.AssertExpectations(h.T())
}

func (h *AuthHandlerSuite) TestLoginUserInvalidInputErr() {
	loginRequest := map[string]string{
		"Email":    "invalid email",
		"Password": "test password",
	}
	body, bodyErr := json.Marshal(loginRequest)
	h.Require().NoError(bodyErr)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (h *AuthHandlerSuite) TestLoginUserErr() {
	tests := []struct {
		name               string
		authServiceError   error
		expectedStatusCode int
	}{
		{
			"invalid credentials",
			domain.ErrInvalidCredentials,
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
				loginRequest := LoginRequest{
					Email:    "testemail@google.com",
					Password: "test password",
				}
				h.mockAuthService.On(
					"Login",
					mock.Anything,
					usecase.LoginUserCommand{
						Email:    loginRequest.Email,
						Password: loginRequest.Password,
					},
				).Return(
					"",
					tt.authServiceError,
				).Once()

				body, bodyErr := json.Marshal(loginRequest)
				h.Require().NoError(bodyErr)

				req := httptest.NewRequest(
					http.MethodPost,
					"/api/v1/auth/login",
					bytes.NewReader(body),
				)
				req.Header.Set("Content-Type", "usecasejson")
				rec := httptest.NewRecorder()
				h.router.ServeHTTP(rec, req)

				h.Assert().Equal(tt.expectedStatusCode, rec.Code)
			},
		)
	}
}
