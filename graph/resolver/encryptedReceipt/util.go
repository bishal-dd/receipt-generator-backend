package encryptedReceipt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)

func convertEdges(edges []*paginationUtil.Edge[*model.EncryptedReceipt]) []*model.EncryptedReceiptEdge {
	dataEdges := make([]*model.EncryptedReceiptEdge, len(edges))
	for i, edge := range edges {
		dataEdges[i] = &model.EncryptedReceiptEdge{
			Cursor: edge.Cursor,
			Node:   edge.Node,
		}
	}
	return dataEdges
}

func searchUpdateInput(input model.UpdateEncryptedReceipt) map[string]interface{} {
	updateInput := make(map[string]interface{})

	// Only add non-nil fields to the updateInput map
	if input.RecipientName != nil {
		updateInput["recipient_name"] = input.RecipientName
	}
	if input.RecipientEmail != nil {
		updateInput["recipient_email"] = input.RecipientEmail
	}
	if input.RecipientAddress != nil {
		updateInput["recipient_address"] = input.RecipientAddress
	}
	if input.RecipientPhone != nil {
		updateInput["recipient_phone"] = input.RecipientPhone
	}
	if input.ReceiptNo != nil {
		updateInput["receipt_no"] = input.ReceiptNo
	}
	if input.PaymentMethod != nil {
		updateInput["payment_method"] = input.PaymentMethod
	}
	if input.PaymentNote != nil {
		updateInput["payment_note"] = input.PaymentNote
	}
	if input.UserID != nil {
		updateInput["user_id"] = input.UserID
	}
	if input.TotalAmount != nil {
		updateInput["total_amount"] = input.TotalAmount
	}

	return updateInput
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

// This is the recommended fix as OAEP is more secure than PKCS#1 v1.5
func encryptRSAWithBase64(publicPEM string, key []byte) []byte {
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
func strPtr(s string) *string {
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
