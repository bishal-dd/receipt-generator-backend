package cloudFront

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
)

func loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKey := os.Getenv("CLOUDFRONT_PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("CLOUDFRONT_PRIVATE_KEY is not set in the environment")
	}

	privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of RSA type")
	}

	return rsaPrivKey, nil
}

func GetCloudFrontURL(key string) (string, error) {
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
    keyID := os.Getenv("CLOUDFRONT_ACCESS_ID")
    privKey, err := loadPrivateKey()
    if err != nil {
        return "", fmt.Errorf("failed to load private key: %s", err)
    }
    signer := sign.NewURLSigner(keyID, privKey)
    baseURL := fmt.Sprintf("https://%s/%s",cloudFrontDomain, key)
    signedURL, err := signer.Sign(baseURL, time.Now().Add(1*time.Hour))
    if err != nil {
        return "", fmt.Errorf("failed to sign URL: %s", err)
    }
    return signedURL, nil
}