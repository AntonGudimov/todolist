package jwt

import (
	"fmt"
	"time"
	"todolist/internal/configs"
	"todolist/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	keystore domain.KeyStore
	config   *configs.JWTConfig
}

func NewJWTService(ks domain.KeyStore, config *configs.JWTConfig) *Service {
	return &Service{keystore: ks, config: config}
}

func (s *Service) CreateAccessToken(userID uuid.UUID) (string, error) {
	activeKey, activeKeyErr := s.keystore.GetSigningKey()
	if activeKeyErr != nil {
		return "", fmt.Errorf("get active private key: %w", activeKeyErr)
	}
	utcNow := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:  userID.String(),
		Issuer:   s.config.Issuer,
		IssuedAt: jwt.NewNumericDate(utcNow),
		ExpiresAt: jwt.NewNumericDate(
			utcNow.Add(s.config.TokenTTL),
		),
	}
	signingMethod, signingMethodErr := getSigningMethod(activeKey.Algorithm)
	if signingMethodErr != nil {
		return "", signingMethodErr
	}
	token := jwt.NewWithClaims(
		signingMethod,
		claims,
	)
	token.Header["kid"] = activeKey.ID
	return token.SignedString(activeKey.PrivateKey)
}

func getSigningMethod(keyAlgorithm domain.Algorithm) (jwt.SigningMethod, error) {
	switch keyAlgorithm {
	case domain.AlgorithmRS256:
		return jwt.SigningMethodRS256, nil
	case domain.AlgorithmRS384:
		return jwt.SigningMethodRS384, nil
	case domain.AlgorithmRS512:
		return jwt.SigningMethodRS512, nil
	case domain.AlgorithmES256:
		return jwt.SigningMethodES256, nil
	case domain.AlgorithmES384:
		return jwt.SigningMethodES384, nil
	case domain.AlgorithmES521:
		return jwt.SigningMethodES512, nil
	default:
		return nil, fmt.Errorf("unsupported key algorithm: %s: %w", keyAlgorithm, domain.ErrNoSigningMethodFound)
	}
}
