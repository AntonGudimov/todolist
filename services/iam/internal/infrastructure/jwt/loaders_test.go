package jwt

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"todolist/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadKeyFromDirOk(t *testing.T) {
	tests := []struct {
		name         string
		keyID        string
		keyAlgorithm domain.Algorithm
		keyPair      KeyPair
	}{
		{
			name:         "rsa 256",
			keyID:        "rsa-2026-08",
			keyAlgorithm: domain.AlgorithmRS256,
			keyPair:      GenerateRsaKey(t, 2048),
		},
		{
			name:         "rsa 384",
			keyID:        "rsa-2026-08",
			keyAlgorithm: domain.AlgorithmRS384,
			keyPair:      GenerateRsaKey(t, 3072),
		},
		{
			name:         "rsa 512",
			keyID:        "rsa-2026-08",
			keyAlgorithm: domain.AlgorithmRS512,
			keyPair:      GenerateRsaKey(t, 4096),
		},
		{
			name:         "ecdsa 224",
			keyID:        "ecdsa-2026-08",
			keyAlgorithm: domain.AlgorithmES224,
			keyPair:      generateEcdsaKey(t, elliptic.P224()),
		},
		{
			name:         "ecdsa 256",
			keyID:        "ecdsa-2026-08",
			keyAlgorithm: domain.AlgorithmES256,
			keyPair:      generateEcdsaKey(t, elliptic.P256()),
		},
		{
			name:         "ecdsa 384",
			keyID:        "ecdsa-2026-08",
			keyAlgorithm: domain.AlgorithmES384,
			keyPair:      generateEcdsaKey(t, elliptic.P384()),
		},
		{
			name:         "ecdsa 521",
			keyID:        "ecdsa-2026-08",
			keyAlgorithm: domain.AlgorithmES521,
			keyPair:      generateEcdsaKey(t, elliptic.P521()),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				dir := filepath.Join(t.TempDir(), tt.keyID)
				dirErr := os.Mkdir(dir, 0755)
				require.NoError(t, dirErr)

				writePrivateKey(t, dir, tt.keyPair.PrivateKey)
				writePublicKey(t, dir, tt.keyPair.PublicKey)
				writeKeyMetadata(t, dir, tt.keyAlgorithm)

				actualKey, keyErr := loadKeyFromDir(dir)
				require.NoError(t, keyErr)

				assert.Equal(t, tt.keyID, actualKey.ID)
				assert.Equal(t, tt.keyAlgorithm, actualKey.Algorithm)

				assertPrivateKeyEqual(t, tt.keyPair.PrivateKey, actualKey.PrivateKey)
				assertPublicKeyEqual(t, tt.keyPair.PublicKey, actualKey.PublicKey)
			},
		)
	}
}

func TestLoadKeyFromDirPrivateKeyMissingErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyAlgorithm := domain.AlgorithmRS512
	keyBits := 4096
	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writePublicKey(t, dir, keyPair.PublicKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrPrivateKeyNotFound)
}

func TestLoadKeyFromDirPublicKeyMissingErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyAlgorithm := domain.AlgorithmRS512
	keyBits := 4096
	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrPublicKeyNotFound)
}

func TestLoadKeyFromDirkKeyMetadataMissingErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writePublicKey(t, dir, keyPair.PublicKey)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrKeyMetadataNotFound)
}

func TestLoadKeyFromDirPrivateKeyPemErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	keyAlgorithm := domain.AlgorithmRS512

	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	var randomBytes []byte
	_, randErr := rand.Read(randomBytes)
	require.NoError(t, randErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writeErr := os.WriteFile(
		filepath.Join(dir, "private.key"),
		randomBytes,
		0644,
	)
	require.NoError(t, writeErr)
	writePublicKey(t, dir, keyPair.PublicKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrInvalidPem)
}

func TestLoadKeyFromDirPublicKeyPemErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	keyAlgorithm := domain.AlgorithmRS512

	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	var randomBytes []byte
	_, randErr := rand.Read(randomBytes)
	require.NoError(t, randErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writeErr := os.WriteFile(
		filepath.Join(dir, "public.key"),
		randomBytes,
		0644,
	)
	require.NoError(t, writeErr)
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrInvalidPem)
}

func TestLoadKeyFromDirPrivateKeyPKCS8Err(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	keyAlgorithm := domain.AlgorithmRS512

	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	var randomBytes []byte
	_, randErr := rand.Read(randomBytes)
	require.NoError(t, randErr)

	keyPair := GenerateRsaKey(t, keyBits)
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: randomBytes,
		},
	)
	writeFileErr := os.WriteFile(
		filepath.Join(dir, "private.key"),
		pemData,
		0600,
	)
	require.NoError(t, writeFileErr)
	writePublicKey(t, dir, keyPair.PublicKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrKeyNotPKCS8PFormat)
}

func TestLoadKeyFromDirPublicKeyPKCS8Err(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	keyAlgorithm := domain.AlgorithmRS512

	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	var randomBytes []byte
	_, randErr := rand.Read(randomBytes)
	require.NoError(t, randErr)

	keyPair := GenerateRsaKey(t, keyBits)
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: randomBytes,
		},
	)
	writeFileErr := os.WriteFile(
		filepath.Join(dir, "public.key"),
		pemData,
		0600,
	)
	require.NoError(t, writeFileErr)
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writeKeyMetadata(t, dir, keyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrKeyNotPKCS8PFormat)
}

func TestLoadKeyFromDirKeyMetadataParseErr(t *testing.T) {
	keyID := "rsa-2026-08"
	keyBits := 4096
	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	keyPair := GenerateRsaKey(t, keyBits)
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writePublicKey(t, dir, keyPair.PublicKey)
	writeFileErr := os.WriteFile(
		filepath.Join(dir, "metadata.json"),
		[]byte("invalid key metadata"),
		0644,
	)
	require.NoError(t, writeFileErr)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrKeyMetadataInvalidFormat)
}

func TestLoadKeyFromDirPrivateKeyNotRSAErr(t *testing.T) {
	keyID := "ecdsa-2026-08"
	wrongKeyAlgorithm := domain.AlgorithmRS512
	dir := filepath.Join(t.TempDir(), keyID)
	dirErr := os.Mkdir(dir, 0755)
	require.NoError(t, dirErr)

	keyPair := generateEcdsaKey(t, elliptic.P384())
	writePrivateKey(t, dir, keyPair.PrivateKey)
	writePublicKey(t, dir, keyPair.PublicKey)
	writeKeyMetadata(t, dir, wrongKeyAlgorithm)

	key, keyErr := loadKeyFromDir(dir)
	assert.Equal(t, domain.Key{}, key)
	assert.ErrorIs(t, keyErr, domain.ErrPrivateKeyNotRSA)
}
