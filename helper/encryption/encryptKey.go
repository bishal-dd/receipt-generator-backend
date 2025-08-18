package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func EncryptKey(publicPEM string, key []byte) []byte {
	block, _ := pem.Decode([]byte(publicPEM))
	if block == nil {
		panic("failed to decode public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	// Use OAEP for encryption to match the decryption function
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, key, nil)
	if err != nil {
		panic(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(encryptedKey))
}
