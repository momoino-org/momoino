package testutils

import (
	"time"
	"wano-island/common/core"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetPlainJwtRSA returns a plain RSA public and private key pair as strings.
// This function is intended for testing purposes only and provides a hardcoded
// public and private key. The returned keys can be used for signing and
// verifying JWT tokens in test scenarios.
//
// Note: The keys are not secure for production use and should not be
// utilized outside of a testing environment.
func GetPlainJwtRSA() (string, string) {
	//nolint:lll // no need to fix
	plainPublicKey := "-----BEGIN PUBLIC KEY-----\nMIIBITANBgkqhkiG9w0BAQEFAAOCAQ4AMIIBCQKCAQB9PuAysl5Hto3FtgwoKJ1p\nM6SIi5BIgzaO0NckpUOC2ikDyazkKg5ST+7MHjb1ophoNDoK/e9mdO21Of4DHvRQ\ndmMxqzSoL+Wt5Pz4HeVIVcUNxcBAiHgagpxijvSgvw+iuPrlZjT1nemcXiiN0miF\nIfsLm7aPBuQZI2ztF1EwbRjAaly8FpzrEo9w0RSLsoUe7GUjG8DlhhJY1H6iKJRl\no550Nc87YafcbIvw2Krkh7Zq8/M4CPy1ApgQ398TFPIkvzLUWggl3cGTWihCvHpx\nIhV+EO6sFKfIsv9TGfD5SAlC5/XMPxc0xk8ZrGk7MFWy0wrxsGur6SzEzo+ChAG3\nAgMBAAE=\n-----END PUBLIC KEY-----"
	//nolint:lll,gosec // no need to fix
	plainPrivateKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQB9PuAysl5Hto3FtgwoKJ1pM6SIi5BIgzaO0NckpUOC2ikDyazk\nKg5ST+7MHjb1ophoNDoK/e9mdO21Of4DHvRQdmMxqzSoL+Wt5Pz4HeVIVcUNxcBA\niHgagpxijvSgvw+iuPrlZjT1nemcXiiN0miFIfsLm7aPBuQZI2ztF1EwbRjAaly8\nFpzrEo9w0RSLsoUe7GUjG8DlhhJY1H6iKJRlo550Nc87YafcbIvw2Krkh7Zq8/M4\nCPy1ApgQ398TFPIkvzLUWggl3cGTWihCvHpxIhV+EO6sFKfIsv9TGfD5SAlC5/XM\nPxc0xk8ZrGk7MFWy0wrxsGur6SzEzo+ChAG3AgMBAAECggEAThp4q4TKAISaMpJN\nUHnLeABpeXE3H9Ebo8IeeE7LI/2yGBebRonndnM8YnPxKAsmac0v6QzkTwtZ9Wrd\nucqC5u58+0tdwghkfaXQD1ZtVkeOZAFO269+3mFW7qthaCDdukcKxyWOnyDDvuyE\n87Qj0+oV6O5I0Tal5ftOgSfKTJzF5G2gSONLyGDtYJp6dXj96KyejlcT/Jf1UCU0\n3oYLV8M+OkvXqNEvUfDO573BZUm79EkrFyerx6r199ibAhBfXjoPS0gg/icpfvU8\noVIwgW3zFwUzWrI3ANpWv8lY/c5nDLPnteJBRiKtb+txOB+Yp6r4rTf8VK4eIKQj\nqJsxUQKBgQDlQ6YwVWlBsO/pefSPg4n9AbmAOtJW4uzrecHQNxKNZx5Yi+x811Cz\nzp2WIwXu+rnIygBqqMJAcO1Qqaxqh/x/2vWhgoKunr92XvaBHdtH2F0hRMAs0rNJ\nrPvoGuGtyXCuyL/Ee5J7ePvxOnRYgU+6JdMmIei1VERpZlUF3j4cuQKBgQCL2ea2\n59RISxZaj+1oIwCtoxL9l/3ujoJqFry1f4tmw7LhdkuTvpSuwWTk7PHG8zZFYtTO\ne1PFYlOx4nilusyRpEK9ZyrNWiR3Rtz8z+kpALjy9wwz2/0j/CYwpOEbkqog/H2X\noAAuNnnSWIIKk+AWz3saVerxgDmGr4wFRjI57wKBgEHy58dXimOfJiQfNL5jtDnX\nWSkNwpvDwyspZxsh/HA4jX4jIe/3b/uJH8OkZ3yLGw4rLVuBF/+5fEqLxFEQtQ2C\nSs4e4MCiYakHQBl8ISvVjVSYlj9OjMxQulXWBb0cCRH+JUu70DM8ZhzKF9WtvOVr\nJAYAExS2HMfE7Ag4Gd3hAoGAEHEg39Ynrgwt553OQpUC6mcmv2vULezRRlm/+/Kv\n1/ggGsPGpOseHeGc1BFLZ6GGeufgrxnuwmEKB/rhRlLM5D6UniH39UaozOEm8A4d\nknWESQRkieBORaHKd6Oa15wJpnEo7t+fxc8fyWwgdc/m46enCHSbd6MkoEIZSzFy\njD8CgYEAnXF3+AyOAI3sTPKvDwTBybR02oHEeflTvCSGZkeJdcp5M3qefGG/0RQo\nVC9jtH7kY3uBMvdDu7V0YrGxDKoevCz9cTjqj+yj5JgliWSErtHSmm0G8nwM0crr\nJHL2ZhWaDr+JDyF+yzFUcW7Ja48wQ8m2Ieea1+OIHi+JhXfm/8E=\n-----END RSA PRIVATE KEY-----"

	return plainPublicKey, plainPrivateKey
}

// GetJWTConfig generates a JWT configuration object with RSA public and private keys.
func GetJWTConfig() *core.JWTConfig {
	plainPublicKey, plainPrivateKey := GetPlainJwtRSA()
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(plainPublicKey))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(plainPrivateKey))

	return &core.JWTConfig{
		PublicKey:             publicKey,
		PrivateKey:            privateKey,
		AccessTokenExpiresIn:  time.Minute,
		RefreshTokenExpiresIn: time.Hour,
	}
}

// GenerateJWT generates a JSON Web Token (JWT) for testing purposes.
// It uses the provided options to customize the JWT claims.
func GenerateJWT(opts ...func(*core.JWTCustomClaims)) string {
	jwtConfig := GetJWTConfig()

	now := time.Now()

	claims := core.JWTCustomClaims{
		Email:             "testing@example.com",
		PreferredUsername: "testing",
		GivenName:         "given_name",
		FamilyName:        "family_name",
		Locale:            "en",
		Roles:             []string{},
		Permissions:       []string{},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.Nil.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.AccessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	for _, fn := range opts {
		fn(&claims)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(jwtConfig.PrivateKey)

	return tokenString
}
