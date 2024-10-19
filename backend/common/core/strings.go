package core

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// RandomString generates a random string of the specified length.
//
// Parameters:
//   - length: The desired length of the random string.
//
// Returns:
//   - A pointer to the generated random string.
//   - An error if there was an issue generating the random string.
func RandomString(length int) (*string, error) {
	buffer := make([]byte, length)

	_, err := rand.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random string with fixed length (%d): %w", length, err)
	}

	randomString := base64.URLEncoding.EncodeToString(buffer)[:length]

	return &randomString, nil
}
