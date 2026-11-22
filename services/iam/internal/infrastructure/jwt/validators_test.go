package jwt

import (
	"crypto/elliptic"
	"os"
	"path/filepath"
	"testing"
	"todolist/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateKeyOk(t *testing.T) {
	tests := []struct {
		name         string
		keyID        string
		keyAlgorithm domain.Algorithm
		keyPair      KeyPair
	}{
		{
			"private key is RSA",
			"rsa-2026-08",
			domain.AlgorithmRS512,
			GenerateRsaKey(t, 4096),
		},
		{
			"private key is ECDSA",
			"ecdsa-2026-08",
			domain.AlgorithmES384,
			generateEcdsaKey(t, elliptic.P384()),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				key := domain.Key{
					ID:         tt.keyID,
					Algorithm:  tt.keyAlgorithm,
					PrivateKey: tt.keyPair.PrivateKey,
					PublicKey:  tt.keyPair.PublicKey,
				}
				validationErr := validateKey(&key)

				assert.NoError(t, validationErr)
			},
		)
	}
}

func TestValidateKeyPrivateKeyParseErr(t *testing.T) {
	tests := []struct {
		name         string
		keyID        string
		keyAlgorithm domain.Algorithm
		keyPair      KeyPair
		expectedErr  error
	}{
		{
			"private key is not RSA",
			"rsa-2026-08",
			domain.AlgorithmRS512,
			generateEcdsaKey(t, elliptic.P384()),
			domain.ErrPrivateKeyNotRSA,
		},
		{
			"private key is not ECDSA",
			"ecdsa-2026-08",
			domain.AlgorithmES384,
			GenerateRsaKey(t, 4096),
			domain.ErrPrivateKeyNotECDSA,
		},
		{
			"unsupported format for private key",
			"unsupported-2026-08",
			"unsupported algorithm",
			KeyPair{},
			domain.ErrPrivateKeyUnsupportedAlgorithm,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				key := domain.Key{
					ID:         tt.keyID,
					Algorithm:  domain.Algorithm(tt.keyAlgorithm),
					PrivateKey: tt.keyPair.PrivateKey,
					PublicKey:  tt.keyPair.PublicKey,
				}
				validationErr := validateKey(&key)
				assert.ErrorIs(t, validationErr, tt.expectedErr)
			},
		)
	}
}

func TestValidateKeyPrivatePublicKeyDontMatchErr(t *testing.T) {
	tests := []struct {
		name string
		key  domain.Key
	}{
		{
			"RSA private and public keys don't match",
			domain.Key{
				ID:         "rsa-2026-08",
				Algorithm:  domain.AlgorithmRS512,
				PrivateKey: GenerateRsaKey(t, 4096).PrivateKey,
				PublicKey:  GenerateRsaKey(t, 4096).PublicKey,
			},
		},
		{
			"ECDSA private and public keys don't match",
			domain.Key{
				ID:         "ecdsa-2026-08",
				Algorithm:  domain.AlgorithmES384,
				PrivateKey: generateEcdsaKey(t, elliptic.P384()).PrivateKey,
				PublicKey:  generateEcdsaKey(t, elliptic.P384()).PublicKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				validationErr := validateKey(&tt.key)
				assert.ErrorIs(t, validationErr, domain.ErrPrivateAndPublicKeysDontMatch)
			},
		)
	}
}

func TestValidateActiveKeyIDOk(t *testing.T) {
	activeKeyID := "rsa-2026-08"
	tmpDir := t.TempDir()
	dir := filepath.Join(tmpDir, activeKeyID)
	dirErr := os.MkdirAll(dir, 0755)
	require.NoError(t, dirErr)

	fileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(activeKeyID),
		0744,
	)
	require.NoError(t, fileErr)

	validateErr := validateActiveKeyID(activeKeyID, tmpDir)
	require.NoError(t, validateErr)
}

func TestValidateActiveKeyIDNotFound(t *testing.T) {
	activeKeyID := "rsa-2026-08"
	tmpDir := t.TempDir()

	fileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(activeKeyID),
		0744,
	)
	require.NoError(t, fileErr)

	validateErr := validateActiveKeyID(activeKeyID, tmpDir)
	assert.ErrorIs(t, validateErr, domain.ErrActiveKeyIDFileNotValid)
}
