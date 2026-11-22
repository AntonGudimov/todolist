package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"todolist/internal/domain"
)

func validateKey(key *domain.Key) error {
	switch key.Algorithm {
	case domain.AlgorithmRS256, domain.AlgorithmRS384, domain.AlgorithmRS512:
		rsaPrivateKey, ok := key.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("algorithm %s requires RSA key: %w", key.Algorithm, domain.ErrPrivateKeyNotRSA)
		}
		rsaValidationErr := validateRsaKey(rsaPrivateKey, key)
		if rsaValidationErr != nil {
			return rsaValidationErr
		}
	case domain.AlgorithmES224, domain.AlgorithmES256, domain.AlgorithmES384, domain.AlgorithmES521:
		ecdsaPrivateKey, ok := key.PrivateKey.(*ecdsa.PrivateKey)
		if !ok {
			return fmt.Errorf("algorithm %s requires ES key: %w", key.Algorithm, domain.ErrPrivateKeyNotECDSA)
		}
		ecdsaValidationErr := validateEcdsaKey(ecdsaPrivateKey, key)
		if ecdsaValidationErr != nil {
			return ecdsaValidationErr
		}
	default:
		return fmt.Errorf("unsupported algorithm %s: %w", key.Algorithm, domain.ErrPrivateKeyUnsupportedAlgorithm)
	}
	return nil
}

func validateRsaKey(rsaKey *rsa.PrivateKey, key *domain.Key) error {
	keysEqual := rsaKey.PublicKey.Equal(key.PublicKey)
	if !keysEqual {
		return domain.ErrPrivateAndPublicKeysDontMatch
	}
	return nil
}

func validateEcdsaKey(ecdsaKey *ecdsa.PrivateKey, key *domain.Key) error {
	keysEqual := ecdsaKey.PublicKey.Equal(key.PublicKey)
	if !keysEqual {
		return domain.ErrPrivateAndPublicKeysDontMatch
	}
	return nil
}

func validateActiveKeyID(activeKeyID string, dir string) error {
	if stat, err := os.Stat(filepath.Join(dir, activeKeyID)); err == nil && stat.IsDir() {
		return nil
	}
	return fmt.Errorf("validate active key ID %s: %w", activeKeyID, domain.ErrActiveKeyIDFileNotValid)
}
