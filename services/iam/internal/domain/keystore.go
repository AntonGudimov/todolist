package domain

import (
	"crypto"
)

type Algorithm string

const (
	AlgorithmRS256 Algorithm = "RS256"
	AlgorithmRS384 Algorithm = "RS384"
	AlgorithmRS512 Algorithm = "RS512"
	AlgorithmES224 Algorithm = "ES224"
	AlgorithmES256 Algorithm = "ES256"
	AlgorithmES384 Algorithm = "ES384"
	AlgorithmES521 Algorithm = "ES521"
)

type Key struct {
	ID         string
	Algorithm  Algorithm
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

type KeyStore interface {
	GetSigningKey() (Key, error)
	Reload() error
}
