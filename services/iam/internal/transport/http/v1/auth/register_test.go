package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func (h *AuthHandlerSuite) TestRegisterOk() {
	registerRequest := RegisterRequest{
		Name:     "test name",
		Email:    "testEmail@google.com",
		Password: "test password",
	}
	registerResult := usecase.CreateUserResult{
		UserID: uuid.New(),
	}
	h.mockAuthService.On(
		"Register",
		mock.Anything,
		usecase.CreateUserCommand{
			Name:     registerRequest.Name,
			Email:    registerRequest.Email,
			Password: registerRequest.Password,
		},
	).Return(
		registerResult,
		nil,
	).Once()

	body, err := json.Marshal(registerRequest)
	h.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusCreated, rec.Code)
	h.mockAuthService.AssertExpectations(h.T())
}

func (h *AuthHandlerSuite) TestRegisterInvalidEmailErr() {
	registerRequest := RegisterRequest{
		Name:     "test name",
		Email:    "testEmail invalid email",
		Password: "test password",
	}

	body, err := json.Marshal(registerRequest)
	h.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (h *AuthHandlerSuite) TestRegisterUserExistsErr() {
	registerRequest := RegisterRequest{
		Name:     "test name",
		Email:    "testEmail@google.com",
		Password: "test password",
	}
	registerResult := usecase.CreateUserResult{}
	h.mockAuthService.On(
		"Register",
		mock.Anything,
		usecase.CreateUserCommand{
			Name:     registerRequest.Name,
			Email:    registerRequest.Email,
			Password: registerRequest.Password,
		},
	).Return(
		registerResult,
		domain.ErrUserAlreadyExists,
	).Once()

	body, err := json.Marshal(registerRequest)
	h.Require().NoError(err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "usecasejson")
	rec := httptest.NewRecorder()

	h.router.ServeHTTP(rec, req)

	h.Assert().Equal(http.StatusConflict, rec.Code)
	h.mockAuthService.AssertExpectations(h.T())
}
