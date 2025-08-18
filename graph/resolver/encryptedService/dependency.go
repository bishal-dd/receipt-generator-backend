package encryptedservice

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type EncryptedServiceResolver struct {
	db            *gorm.DB
	httpClient    *resty.Client
	publicKeyPEM  string
	privateKeyPEM string
}

func InitializeEncryptedServiceResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *EncryptedServiceResolver {
	return &EncryptedServiceResolver{
		db:            db,
		httpClient:    httpClient,
		publicKeyPEM:  publicKeyPEM,
		privateKeyPEM: privateKeyPEM,
	}
}

const EncryptedServicesKey = "services"
const EncryptedServiceKey = "service"
