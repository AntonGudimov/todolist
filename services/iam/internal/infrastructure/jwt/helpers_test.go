package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"todolist/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func GenerateRsaKey(t *testing.T, bits int) KeyPair {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, bits)
	require.NoError(t, err)
	return KeyPair{PrivateKey: key, PublicKey: &key.PublicKey}
}

func generateEcdsaKey(t *testing.T, c elliptic.Curve) KeyPair {
	t.Helper()
	key, err := ecdsa.GenerateKey(c, rand.Reader)
	require.NoError(t, err)
	return KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}
}

func writePrivateKey(t *testing.T, dir string, privateKey any) {
	t.Helper()
	der, derErr := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, derErr)
	data := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: der,
		},
	)
	err := os.WriteFile(
		filepath.Join(
			dir, "private.key",
		),
		data,
		0600,
	)
	require.NoError(t, err)
}

func writePublicKey(t *testing.T, dir string, publicKey any) {
	t.Helper()
	der, derErr := x509.MarshalPKIXPublicKey(publicKey)
	require.NoError(t, derErr)
	data := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: der,
		},
	)
	err := os.WriteFile(
		filepath.Join(
			dir, "public.key",
		),
		data,
		0600,
	)
	require.NoError(t, err)
}

func writeKeyMetadata(t *testing.T, dir string, algorithm domain.Algorithm) {
	t.Helper()

	keyMetadata := KeyMetaData{
		Algorithm: domain.Algorithm(algorithm),
	}
	data, err := json.Marshal(keyMetadata)
	require.NoError(t, err)
	os.WriteFile(
		filepath.Join(dir, "metadata.json"),
		data,
		0600,
	)
}

func assertPrivateKeyEqual(t *testing.T, expectedKey crypto.PrivateKey, actualKey crypto.PrivateKey) {
	t.Helper()

	expectedBytes, expectedBytesErr := x509.MarshalPKCS8PrivateKey(expectedKey)
	require.NoError(t, expectedBytesErr)

	actualBytes, actualBytesErr := x509.MarshalPKCS8PrivateKey(actualKey)
	require.NoError(t, actualBytesErr)

	assert.Equal(t, expectedBytes, actualBytes)
}

func assertPublicKeyEqual(t *testing.T, expectedKey crypto.PublicKey, actualKey crypto.PublicKey) {
	t.Helper()

	expectedBytes, expectedBytesErr := x509.MarshalPKIXPublicKey(expectedKey)
	require.NoError(t, expectedBytesErr)

	actualBytes, actualBytesErr := x509.MarshalPKIXPublicKey(actualKey)
	require.NoError(t, actualBytesErr)

	assert.Equal(t, expectedBytes, actualBytes)
}
