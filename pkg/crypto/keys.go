package crypto

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// LoadPrivateKey loads RSA private key from env
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	key := os.Getenv("PRIVATE_KEY")
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey.(*rsa.PrivateKey), nil
}

// LoadPublicKey loads RSA public key from env
func LoadPublicKey() (*rsa.PublicKey, error) {
	key := os.Getenv("PUBLIC_KEY")
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func PublicKeyToPEM(pubKey *rsa.PublicKey) (string, error) {
	// Marshal the public key to PKIX, ASN.1 DER form
	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	// Create a PEM block with the DER bytes
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	var pemBuffer bytes.Buffer
	if err := pem.Encode(&pemBuffer, pemBlock); err != nil {
		return "", err
	}

	return pemBuffer.String(), nil
}

func PrivateKeyToPEM(privKey *rsa.PrivateKey) (string, error) {
	// Marshal the private key to PKCS8, ASN.1 DER form
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", err
	}

	// Create a PEM block with the DER bytes
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}

	var pemBuffer bytes.Buffer
	if err := pem.Encode(&pemBuffer, pemBlock); err != nil {
		return "", err
	}

	return pemBuffer.String(), nil
}
