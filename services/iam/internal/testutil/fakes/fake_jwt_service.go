package fakes

import "github.com/google/uuid"

type FakeJWTService struct {
	AccessToken          string
	CreateAccessTokenErr error
}

func NewFakeJWTService(accessToken string, createAccessTokenErr error) *FakeJWTService {
	return &FakeJWTService{
		AccessToken:          accessToken,
		CreateAccessTokenErr: createAccessTokenErr,
	}
}

func (s *FakeJWTService) CreateAccessToken(userID uuid.UUID) (string, error) {
	return s.AccessToken, s.CreateAccessTokenErr
}
