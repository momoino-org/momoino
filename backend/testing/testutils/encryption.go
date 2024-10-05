package testutils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"wano-island/common/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func generateRandomString(length int) string {
	buffer := make([]byte, length)

	_, err := rand.Read(buffer)
	Expect(err).NotTo(HaveOccurred())

	return base64.URLEncoding.EncodeToString(buffer)[:length]
}

// GenerateSecretFiles creates RSA public and private keys, and a random secret key.
// The generated keys and secret key are stored in the current directory as "public.key",
// "private.key", and "secret.key" respectively.
// The keys and secret key are stored with 0600 file permissions, meaning they are only
// readable and writable by the owner.
// The function returns an error if any of the file operations fail.
// After the function execution, it cleans up the generated files by removing them.
func GenerateSecretFiles() {
	publicKey, privateKey, err := core.GenerateRSAKey()
	Expect(err).NotTo(HaveOccurred())

	//nolint:mnd // For testing purposes
	Expect(os.WriteFile("public.key", publicKey, 0600)).NotTo(HaveOccurred())
	//nolint:mnd // For testing purposes
	Expect(os.WriteFile("private.key", privateKey, 0600)).NotTo(HaveOccurred())
	//nolint:mnd // For testing purposes
	Expect(os.WriteFile("secret.key", []byte(generateRandomString(32)), 0600)).NotTo(HaveOccurred())

	DeferCleanup(func() {
		os.Remove("public.key")
		os.Remove("private.key")
		os.Remove("secret.key")
	})
}
