package receiptFiles

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

func parseRSAPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

func jsonEqual(a, b []byte) bool {
	var o1, o2 any
	if err := json.Unmarshal(a, &o1); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &o2); err != nil {
		return false
	}
	return fmt.Sprintf("%v", o1) == fmt.Sprintf("%v", o2)
}
