package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func DecryptField(encryptedBase64 *string, aesKey []byte, iv []byte) *string {
	if encryptedBase64 == nil {
		return nil
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(*encryptedBase64)
	if err != nil {
		return nil
	}
	decrypted, err := decryptAES(aesKey, iv, encryptedBytes)
	if err != nil {
		return nil
	}
	str := string(decrypted)
	return &str
}

func decryptAES(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)
	return pkcs7Unpad(decrypted)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot unpad empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-padding], nil
}
