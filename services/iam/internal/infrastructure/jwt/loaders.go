package jwt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"todolist/internal/domain"
)

type KeyMetaData struct {
	Algorithm domain.Algorithm `json:"algorithm"`
}

func loadActiveKeyID(dir string) (string, error) {
	activeData, activeDataErr := os.ReadFile(
		filepath.Join(dir, "active"),
	)
	if activeDataErr != nil {
		return "", fmt.Errorf("read active key ID: %w", activeDataErr)
	}
	activeKeyID := strings.TrimSpace(string(activeData))
	return activeKeyID, nil
}

func loadPrivateKeyFromDir(dir string) ([]byte, error) {
	keyPath := filepath.Join(dir, "private.key")
	privateKeyPem, privateKeyPemErr := os.ReadFile(
		keyPath,
	)
	if privateKeyPemErr != nil {
		return nil, fmt.Errorf("read private key %s: %w", keyPath, domain.ErrPrivateKeyNotFound)
	}
	return privateKeyPem, nil
}

func loadPublicKeyFromDir(dir string) ([]byte, error) {
	keyPath := filepath.Join(dir, "public.key")
	publicKeyPem, publicKeyPemErr := os.ReadFile(
		keyPath,
	)
	if publicKeyPemErr != nil {
		return nil, fmt.Errorf("read public key %s: %w", keyPath, domain.ErrPublicKeyNotFound)
	}
	return publicKeyPem, nil
}

func loadKeyMetadata(dir string) ([]byte, error) {
	keyMetadataPath := filepath.Join(dir, "metadata.json")
	keyMetadataBytes, keyMetaDataBytesErr := os.ReadFile(
		keyMetadataPath,
	)
	if keyMetaDataBytesErr != nil {
		return nil, fmt.Errorf("read key metadata %s: %w", keyMetadataPath, domain.ErrKeyMetadataNotFound)
	}
	return keyMetadataBytes, nil
}

func loadKeyFromDir(dir string) (domain.Key, error) {
	keyID := filepath.Base(dir)
	privateKeyBytes, privateKeyBytesErr := loadPrivateKeyFromDir(dir)
	if privateKeyBytesErr != nil {
		return domain.Key{}, privateKeyBytesErr
	}
	publicKeyBytes, publicKeyBytesErr := loadPublicKeyFromDir(dir)
	if publicKeyBytesErr != nil {
		return domain.Key{}, publicKeyBytesErr
	}
	keyMetadataBytes, keyMetadataBytesErr := loadKeyMetadata(dir)
	if keyMetadataBytesErr != nil {
		return domain.Key{}, keyMetadataBytesErr
	}
	privateKey, privateKeyErr := parsePrivateKey(privateKeyBytes)
	if privateKeyErr != nil {
		return domain.Key{}, fmt.Errorf("parse private key: keyID %s: %w", keyID, privateKeyErr)
	}
	publicKey, publicKeyErr := parsePublicKey(publicKeyBytes)
	if publicKeyErr != nil {
		return domain.Key{}, fmt.Errorf("parse public: keyID %s: %w", keyID, publicKeyErr)
	}
	keyMetadata, keyMetadataErr := parseKeyMetadata(keyMetadataBytes)
	if keyMetadataErr != nil {
		return domain.Key{}, fmt.Errorf("parse key metadata: keyID: %s: %w", keyID, keyMetadataErr)
	}
	key := domain.Key{
		ID:         keyID,
		Algorithm:  keyMetadata.Algorithm,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
	validateErr := validateKey(&key)
	if validateErr != nil {
		return domain.Key{}, fmt.Errorf("validate key ID %s: %w", keyID, validateErr)
	}
	return key, nil
}
