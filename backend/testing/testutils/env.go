package testutils

import (
	"wano-island/common/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ConfigureMinimumEnvVariables sets up the minimum required environment
// variables for testing purposes. This function generates a random AES
// secret key and RSA key pair, then sets the following environment
// variables using the generated values:
//   - APP_SECRET_KEY: The generated AES secret key.
//   - APP_JWT_RSA_PUBLIC_KEY: The generated RSA public key.
//   - APP_JWT_RSA_PRIVATE_KEY: The generated RSA private key.
//
// This function is intended for use in test setups where secure
// keys and values are required.
func ConfigureMinimumEnvVariables() {
	secretKey, err := core.RandomString(core.AESSecretKeyLength)
	Expect(err).NotTo(HaveOccurred())

	publicKey, privateKey, err := core.GenerateRSAKey()
	Expect(err).NotTo(HaveOccurred())

	t := GinkgoT()
	t.Setenv("APP_SECRET_KEY", *secretKey)
	t.Setenv("APP_JWT_RSA_PUBLIC_KEY", string(publicKey))
	t.Setenv("APP_JWT_RSA_PRIVATE_KEY", string(privateKey))
}
