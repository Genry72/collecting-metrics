package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Encrypt Кодируем
func Encrypt(value []byte, password string) (string, error) {

	aesgcm, nonce, err := getGcmAndNonce(password)
	if err != nil {
		return "", fmt.Errorf("getGcmAndNonce: %w", err)
	}

	return hex.EncodeToString(aesgcm.Seal(nil, nonce, value, nil)), nil
}

// Decrypt Декодируем
func Decrypt(value, password string) ([]byte, error) {
	aesgcm, nonce, err := getGcmAndNonce(password)
	if err != nil {
		return nil, fmt.Errorf("getGcmAndNonce: %w", err)
	}

	encrypted, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString: %w", err)
	}

	return aesgcm.Open(nil, nonce, encrypted, nil)
}

func getGcmAndNonce(password string) (cipher.AEAD, []byte, error) {
	key := sha256Hashing(password)

	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	return aesgcm, nonce, nil
}

func sha256Hashing(input string) []byte {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return sha256Hash[:]

}
