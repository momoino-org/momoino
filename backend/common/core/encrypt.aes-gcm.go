package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/samber/lo"
)

type aesGCMEncryptor struct {
	config AppConfig
}

var _ Encryptor = (*aesGCMEncryptor)(nil)

func newAESEncryptor(config AppConfig) *aesGCMEncryptor {
	return &aesGCMEncryptor{
		config: config,
	}
}

// Encrypt takes a plaintext string and returns its encrypted version as a base64-encoded string.
// It uses AES-GCM encryption with a secret key obtained from the application's configuration.
func (e *aesGCMEncryptor) Encrypt(plaintext []byte) (*string, error) {
	var err error

	block, err := aes.NewCipher(e.config.GetSecretKey())
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return lo.ToPtr(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt takes an encrypted text string (base64-encoded) and returns its decrypted version as a string pointer.
// It uses AES-GCM decryption with a secret key obtained from the application's configuration.
func (e *aesGCMEncryptor) Decrypt(encryptedText string) (*string, error) {
	decodedText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.config.GetSecretKey())
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(decodedText) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherText := decodedText[:nonceSize], decodedText[nonceSize:]

	plainText, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return lo.ToPtr(string(plainText)), nil
}
