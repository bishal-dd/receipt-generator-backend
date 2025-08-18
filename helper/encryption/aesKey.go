package encryption

import (
	"crypto/aes"
	"crypto/rand"
)

func GenerateAESKeyAndIV() ([]byte, []byte, error) {
	aesKey := make([]byte, 32) // AES-256
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, nil, err
	}
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}
	return aesKey, iv, nil
}
