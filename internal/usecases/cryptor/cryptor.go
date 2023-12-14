package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// Encrypt Кодируем байтовый массив
func Encrypt(value []byte, password string) (string, error) {

	aesgcm, nonce, err := getGcmAndNonce(password)
	if err != nil {
		return "", fmt.Errorf("getGcmAndNonce: %w", err)
	}

	return hex.EncodeToString(aesgcm.Seal(nil, nonce, value, nil)), nil
}

// Decrypt Декодируем байтовый массив
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

// DecryptWithPrivateKey Расшифровка при помощи закрытого ключа
func DecryptWithPrivateKey(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return decryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
}

// EncryptBodyWithPublicKey шифрование при помощи открытого ключа
func EncryptBodyWithPublicKey(data []byte, pub *rsa.PublicKey) ([]byte, error) {
	return encryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
}

// https://stackoverflow.com/questions/62348923/rs256-message-too-long-for-rsa-public-key-size-error-signing-jwt
func encryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2

	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, fmt.Errorf("rsa.EncryptOAEP: %w", err)
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

// https://stackoverflow.com/questions/62348923/rs256-message-too-long-for-rsa-public-key-size-error-signing-jwt
func decryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()

	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, fmt.Errorf("rsa.DecryptOAEP: %w", err)
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
