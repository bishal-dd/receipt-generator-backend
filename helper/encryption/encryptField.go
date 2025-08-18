package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptField(value *string, aesKey, iv []byte) *string {
	if value == nil {
		return nil
	}
	encrypted, _ := encryptAES(aesKey, iv, []byte(*value))
	base64Str := base64.StdEncoding.EncodeToString(encrypted)
	return &base64Str
}

func encryptAES(key, iv, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	padLen := aes.BlockSize - len(plaintext)%aes.BlockSize
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	padded := append(plaintext, padding...)

	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)

	return ciphertext, nil
}
