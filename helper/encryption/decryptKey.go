package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func DecryptKeyAndIV(privateKeyPEM string, AesKeyEncrypted string, AesIv string) ([]byte, []byte, error) {
	aesKey, err := decryptRSAWithBase64([]byte(privateKeyPEM), []byte(AesKeyEncrypted))
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt AES key: %w", err)
	}

	// 2. Decode IV
	iv, err := base64.StdEncoding.DecodeString(AesIv)
	if err != nil {
		return nil, nil, fmt.Errorf("decode IV: %w", err)
	}

	return aesKey, iv, nil
}

func decryptRSAWithBase64(privateKeyPEM []byte, encryptedBase64 []byte) ([]byte, error) {
	// Decode PEM
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Parse PKCS#8 private key
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA private key")
	}

	// Decode base64 encrypted message
	encryptedBytes, err := base64.StdEncoding.DecodeString(string(encryptedBase64))
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	// Decrypt using RSA OAEP
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, encryptedBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA decrypt: %w", err)
	}

	return decrypted, nil
}
