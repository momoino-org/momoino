package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"go.uber.org/fx"
)

type Encryptor interface {
	Encrypt(plaintext []byte) (*string, error)
	Decrypt(encryptedText string) (*string, error)
}

type noopEncryptor struct{}

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

func (e noopEncryptor) Encrypt(plaintext []byte) (*string, error) {
	return lo.ToPtr(fmt.Sprintf("encrypted:%v", plaintext)), nil
}

func (e noopEncryptor) Decrypt(encryptedText string) (*string, error) {
	decryptedText := strings.TrimPrefix(encryptedText, "encrypted:")
	return &decryptedText, nil
}

func NewNoopEncryptor() *noopEncryptor {
	return &noopEncryptor{}
}

// NewEncryptionModule is an Fx option that provides an encryption module for the application.
func NewEncryptionModule() fx.Option {
	return fx.Module(
		"Encryption module",
		fx.Provide(
			fx.Annotate(newAESEncryptor, fx.As(new(Encryptor)), fx.ResultTags(`name:"aes-gcm"`)),
		),
	)
}
