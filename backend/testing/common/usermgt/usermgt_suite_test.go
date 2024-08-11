package usermgt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"wano-island/common/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = BeforeSuite(func() {
	IgnoreGinkgoParallelClient()
})

func TestUsermgt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usermgt Suite")
}

// generateRSAKey generates a pair of RSA public and private keys.
// The private key is encoded in PEM format and the public key is also encoded in PEM format.
func generateRSAKey() ([]byte, []byte) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	publicKey := &privateKey.PublicKey
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	return pem.EncodeToMemory(publicKeyPEM), pem.EncodeToMemory(privateKeyPEM)
}

// generateJWTConfig generates a JWT configuration object with RSA public and private keys.
func generateJWTConfig() *core.JWTConfig {
	publicKey, privateKey := generateRSAKey()

	return &core.JWTConfig{
		PublicKey:             publicKey,
		PrivateKey:            privateKey,
		AccessTokenExpiresIn:  3600,
		RefreshTokenExpiresIn: 86400,
	}
}
