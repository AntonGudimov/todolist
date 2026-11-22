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

func TestKeyStoreReloadStartUpOk(t *testing.T) {
	rsaKey := GenerateRsaKey(t, 4096)
	rsaKeyID := "rsa-2026-08"
	rsaAlgorithm := domain.AlgorithmRS512

	ecdsaKey := generateEcdsaKey(t, elliptic.P384())
	ecdsaKeyID := "ecdsa-2026-07"
	ecdsaKeyAlgorithm := domain.AlgorithmES384

	tmpDir := t.TempDir()
	rsaKeyDir := filepath.Join(tmpDir, rsaKeyID)
	rsaKeyDirErr := os.Mkdir(
		rsaKeyDir,
		0755,
	)
	require.NoError(t, rsaKeyDirErr)

	ecdsayKeyDir := filepath.Join(tmpDir, ecdsaKeyID)
	ecdsaKeyDirErr := os.Mkdir(
		ecdsayKeyDir,
		0755,
	)
	require.NoError(t, ecdsaKeyDirErr)

	writePrivateKey(t, rsaKeyDir, rsaKey.PrivateKey)
	writePublicKey(t, rsaKeyDir, rsaKey.PublicKey)
	writeKeyMetadata(t, rsaKeyDir, rsaAlgorithm)

	writePrivateKey(t, ecdsayKeyDir, ecdsaKey.PrivateKey)
	writePublicKey(t, ecdsayKeyDir, ecdsaKey.PublicKey)
	writeKeyMetadata(t, ecdsayKeyDir, ecdsaKeyAlgorithm)

	activeFileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(rsaKeyID),
		0755,
	)
	require.NoError(t, activeFileErr)

	keyStore, keyStoreErr := NewFileKeyStore(tmpDir)
	require.NoError(t, keyStoreErr)

	assert.Len(t, keyStore.keys, 2)
	assert.Equal(t, rsaKeyID, keyStore.signingKeyID)
	assert.Contains(t, keyStore.keys, rsaKeyID)
	assert.Contains(t, keyStore.keys, ecdsaKeyID)
}

func TestKeyStoreReloadAddNewKeyOk(t *testing.T) {
	rsaKey := GenerateRsaKey(t, 4096)
	rsaKeyID := "rsa-2026-08"
	rsaAlgorithm := domain.AlgorithmRS512

	tmpDir := t.TempDir()
	rsaKeyDir := filepath.Join(tmpDir, rsaKeyID)
	rsaKeyDirErr := os.Mkdir(
		rsaKeyDir,
		0755,
	)
	require.NoError(t, rsaKeyDirErr)

	writePrivateKey(t, rsaKeyDir, rsaKey.PrivateKey)
	writePublicKey(t, rsaKeyDir, rsaKey.PublicKey)
	writeKeyMetadata(t, rsaKeyDir, rsaAlgorithm)

	activeFileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(rsaKeyID),
		0755,
	)
	require.NoError(t, activeFileErr)

	keyStore, keyStoreErr := NewFileKeyStore(tmpDir)
	require.NoError(t, keyStoreErr)

	ecdsaKey := generateEcdsaKey(t, elliptic.P384())
	ecdsaKeyID := "ecdsa-2026-07"
	ecdsaKeyAlgorithm := domain.AlgorithmES384

	ecdsayKeyDir := filepath.Join(tmpDir, ecdsaKeyID)
	ecdsaKeyDirErr := os.Mkdir(
		ecdsayKeyDir,
		0755,
	)
	require.NoError(t, ecdsaKeyDirErr)
	writePrivateKey(t, ecdsayKeyDir, ecdsaKey.PrivateKey)
	writePublicKey(t, ecdsayKeyDir, ecdsaKey.PublicKey)
	writeKeyMetadata(t, ecdsayKeyDir, ecdsaKeyAlgorithm)
	activeFileErr = os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(ecdsaKeyID),
		0755,
	)
	require.NoError(t, activeFileErr)

	keyStoreReloadErr := keyStore.Reload()
	require.NoError(t, keyStoreReloadErr)

	assert.Equal(t, ecdsaKeyID, keyStore.signingKeyID)
	assert.Contains(t, keyStore.keys, ecdsaKeyID)
}

func TestKeyStoreReloadRemoveOldKey(t *testing.T) {
	rsaKey := GenerateRsaKey(t, 4096)
	rsaKeyID := "rsa-2026-08"
	rsaAlgorithm := domain.AlgorithmRS512

	ecdsaKey := generateEcdsaKey(t, elliptic.P384())
	ecdsaKeyID := "ecdsa-2026-07"
	ecdsaKeyAlgorithm := domain.AlgorithmES384

	tmpDir := t.TempDir()
	rsaKeyDir := filepath.Join(tmpDir, rsaKeyID)
	rsaKeyDirErr := os.Mkdir(
		rsaKeyDir,
		0755,
	)
	require.NoError(t, rsaKeyDirErr)

	ecdsayKeyDir := filepath.Join(tmpDir, ecdsaKeyID)
	ecdsaKeyDirErr := os.Mkdir(
		ecdsayKeyDir,
		0755,
	)
	require.NoError(t, ecdsaKeyDirErr)

	writePrivateKey(t, rsaKeyDir, rsaKey.PrivateKey)
	writePublicKey(t, rsaKeyDir, rsaKey.PublicKey)
	writeKeyMetadata(t, rsaKeyDir, rsaAlgorithm)

	writePrivateKey(t, ecdsayKeyDir, ecdsaKey.PrivateKey)
	writePublicKey(t, ecdsayKeyDir, ecdsaKey.PublicKey)
	writeKeyMetadata(t, ecdsayKeyDir, ecdsaKeyAlgorithm)

	activeFileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(ecdsaKeyID),
		0755,
	)
	require.NoError(t, activeFileErr)

	keyStore, keyStoreErr := NewFileKeyStore(tmpDir)
	require.NoError(t, keyStoreErr)

	removeRsaKeyDirErr := os.RemoveAll(rsaKeyDir)
	require.NoError(t, removeRsaKeyDirErr)

	err := keyStore.Reload()
	require.NoError(t, err)

	assert.NotContains(t, keyStore.keys, rsaKeyID)
	assert.Len(t, keyStore.keys, 1)
	assert.Equal(t, ecdsaKeyID, keyStore.signingKeyID)
}

func TestKeyStoreReloadInvalidKeyErrKeepsOldState(t *testing.T) {
	rsaKey := GenerateRsaKey(t, 4096)
	rsaKeyID := "rsa-2026-08"
	rsaAlgorithm := domain.AlgorithmRS512

	tmpDir := t.TempDir()
	rsaKeyDir := filepath.Join(tmpDir, rsaKeyID)
	rsaKeyDirErr := os.Mkdir(
		rsaKeyDir,
		0755,
	)
	require.NoError(t, rsaKeyDirErr)

	writePrivateKey(t, rsaKeyDir, rsaKey.PrivateKey)
	writePublicKey(t, rsaKeyDir, rsaKey.PublicKey)
	writeKeyMetadata(t, rsaKeyDir, rsaAlgorithm)

	activeFileErr := os.WriteFile(
		filepath.Join(tmpDir, "active"),
		[]byte(rsaKeyID),
		0755,
	)
	require.NoError(t, activeFileErr)

	keyStore, keyStoreErr := NewFileKeyStore(tmpDir)
	require.NoError(t, keyStoreErr)

	ecdsaKeyID := "ecdsa-2026-07"
	ecdsayKeyDir := filepath.Join(tmpDir, ecdsaKeyID)
	ecdsaKeyDirErr := os.Mkdir(
		ecdsayKeyDir,
		0755,
	)
	require.NoError(t, ecdsaKeyDirErr)

	oldKeys := keyStore.keys
	oldSigningKeyID := keyStore.signingKeyID

	err := keyStore.Reload()
	assert.Error(t, err)

	assert.Equal(t, oldKeys, keyStore.keys)
	assert.Equal(t, oldSigningKeyID, keyStore.signingKeyID)
}

func TestKeyStoreGetActvikeKeyOk(t *testing.T) {
	keyPair := GenerateRsaKey(t, 4096)
	key := domain.Key{
		ID:         "rsa-2026-08",
		Algorithm:  domain.AlgorithmRS512,
		PrivateKey: keyPair.PrivateKey,
		PublicKey:  keyPair.PublicKey,
	}
	keyStore := FileKeyStore{
		signingKeyID: key.ID,
		keys:         map[string]domain.Key{key.ID: key},
	}
	signingKey, err := keyStore.GetSigningKey()
	require.NoError(t, err)
	assert.Equal(t, key, signingKey)
}
