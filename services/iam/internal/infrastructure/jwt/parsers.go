package jwt

import (
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"todolist/internal/domain"
)

func parsePrivateKey(data []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode pem: %w", domain.ErrInvalidPem)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse pkcs8: %w", domain.ErrKeyNotPKCS8PFormat)
	}
	return key, nil
}

func parsePublicKey(data []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode pem: %w", domain.ErrInvalidPem)
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse pkcs8: %w", domain.ErrKeyNotPKCS8PFormat)
	}
	return key, nil
}

func parseKeyMetadata(data []byte) (KeyMetaData, error) {
	var keyMetadata KeyMetaData

	keyMetaDataErr := json.Unmarshal(data, &keyMetadata)
	if keyMetaDataErr != nil {
		return KeyMetaData{}, fmt.Errorf("parse key metadata: %w", domain.ErrKeyMetadataInvalidFormat)
	}
	return keyMetadata, nil
}
