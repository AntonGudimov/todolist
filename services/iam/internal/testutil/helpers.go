package testutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"

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
