package testutils

import (
	"wano-island/common/core"

	. "github.com/onsi/ginkgo/v2"
)

// ConfigureMinimumEnvVariables sets up the minimum required environment
// variables for testing purposes.
func ConfigureMinimumEnvVariables() {
	t := GinkgoT()
	t.Setenv("APP_MODE", core.TestingMode)
	t.Setenv("APP_SECRET_KEY", "eLPgQbF,g!Yz)6E%9Ghj5.KZMWvw$!9y")
	t.Setenv("APP_KEYCLOAK_WELL_KNOWN_URL", "http://keycloak.test/.well-known/openid-configuration")
}
