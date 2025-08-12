package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

var masterPubKey *rsa.PublicKey
var masterPrivKey *rsa.PrivateKey

// Initialize keys (call this in your app startup, loading from PEM/env)
func InitKeys(pub *rsa.PublicKey, priv *rsa.PrivateKey) {
	masterPubKey = pub
	masterPrivKey = priv
}

// Encrypt a single field with envelope encryption
func EncryptField(plaintext string) (encryptedData, encryptedDEK, nonce string, err error) {
	// Generate DEK (32 bytes for AES-256)
	dek := make([]byte, 32)
	if _, err = rand.Read(dek); err != nil {
		return "", "", "", err
	}

	// Encrypt plaintext with AES-GCM
	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", "", err
	}

	nonceBytes := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonceBytes); err != nil {
		return "", "", "", err
	}

	ciphertext := gcm.Seal(nil, nonceBytes, []byte(plaintext), nil)

	// Encrypt DEK with RSA
	encryptedDEKBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, masterPubKey, dek, nil)
	if err != nil {
		return "", "", "", err
	}

	// Base64 encode for storage
	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(encryptedDEKBytes),
		base64.StdEncoding.EncodeToString(nonceBytes), nil
}

// Decrypt a field
func DecryptField(ciphertextB64, encryptedDEKB64, nonceB64 string) (string, error) {
	if masterPrivKey == nil {
		return "", errors.New("private key not initialized")
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(ciphertextB64)
	encryptedDEK, _ := base64.StdEncoding.DecodeString(encryptedDEKB64)
	nonce, _ := base64.StdEncoding.DecodeString(nonceB64)

	// Decrypt DEK
	dek, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, masterPrivKey, encryptedDEK, nil)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
