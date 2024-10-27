package testutils

import (
	"wano-island/common/core"

	. "github.com/onsi/ginkgo/v2"
)

// ConfigureMinimumEnvVariables sets up the minimum required environment
// variables for testing purposes.
func ConfigureMinimumEnvVariables() {
	publicKey, privateKey := GetPlainJwtRSA()

	t := GinkgoT()
	t.Setenv("APP_MODE", core.TestingMode)
	t.Setenv("APP_SECRET_KEY", "eLPgQbF,g!Yz)6E%9Ghj5.KZMWvw$!9y")
	t.Setenv("APP_JWT_RSA_PUBLIC_KEY", publicKey)
	t.Setenv("APP_JWT_RSA_PRIVATE_KEY", privateKey)
	t.Setenv("APP_JWT_ACCESS_TOKEN_EXPIRES_IN", "1h")
	t.Setenv("APP_SESSION_LIFETIME", "336h")
	t.Setenv("APP_SESSION_IDLE_TIMEOUT", "168h")
}
