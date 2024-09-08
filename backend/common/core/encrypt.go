package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// GenerateRSAKey generates a pair of RSA public and private keys.
// The private key is encoded in PEM format and the public key is also encoded in PEM format.
func GenerateRSAKey() ([]byte, []byte, error) {
	//nolint:mnd // No need to fix
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate RSA private key: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKey := &privateKey.PublicKey
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	return pem.EncodeToMemory(publicKeyPEM), pem.EncodeToMemory(privateKeyPEM), nil
}
