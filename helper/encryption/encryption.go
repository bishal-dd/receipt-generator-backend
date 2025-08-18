package encryption

import (
	"crypto/rsa"
)

var masterPubKey *rsa.PublicKey
var masterPrivKey *rsa.PrivateKey

// Initialize keys (call this in your app startup, loading from PEM/env)
func InitKeys(pub *rsa.PublicKey, priv *rsa.PrivateKey) {
	masterPubKey = pub
	masterPrivKey = priv
}
