package receiptPDFGenerator

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
	"github.com/clerk/clerk-sdk-go/v2"
)

func calculateTotalAmount(services []*model.CreateBulkService, tax float64) (float64, float64, float64) {
	subtotal := 0.0
	for _, serviceInput := range services {
		subtotal += serviceInput.Amount
	}

	taxRate := tax / 100
	taxAmount := subtotal * taxRate
	totalAmount := subtotal + taxAmount

	return totalAmount, subtotal, taxAmount
}

func updateProfileImages(profile *model.Profile, organization *clerk.Organization) error {
	if profile.SignatureImage != nil && *profile.SignatureImage != "" {
		signedURL, err := cloudFront.GetCloudFrontURL(*profile.SignatureImage)
		if err != nil {
			return err
		}
		profile.SignatureImage = &signedURL
	}
	profile.CompanyName = &organization.Name
	if organization.HasImage {
		profile.LogoImage = organization.ImageURL
	}

	return nil
}

func emailInputToReceiptModel(input model.SendReceiptPDFToEmail, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	return &model.Receipt{
		ID:               ids.UUID(),
		UserID:           userId,
		CreatedAt:        time.Now().Format(time.RFC3339),
		ReceiptName:      receiptInput.ReceiptName,
		RecipientPhone:   receiptInput.RecipientPhone,
		RecipientName:    receiptInput.RecipientName,
		RecipientEmail:   &receiptInput.RecipientEmail,
		RecipientAddress: receiptInput.RecipientAddress,
		ReceiptNo:        receiptInput.ReceiptNo,
		Date:             receiptInput.Date,
		PaymentMethod:    receiptInput.PaymentMethod,
		PaymentNote:      receiptInput.PaymentNote,
		TotalAmount:      &totalAmount,
		SubTotalAmount:   &subtotal,
		TaxAmount:        &taxAmount,
		Services:         make([]*model.Service, 0),
	}
}

func emailInputToEncryptedReceiptModel(input model.SendReceiptPDFToEmail, userId string, totalAmount, subtotal, taxAmount float64, publicKeyPEM string) (*model.EncryptedReceipt, error) {
	receiptInput := input

	aesKey, iv, err := encryption.GenerateAESKeyAndIV()
	if err != nil {
		return nil, err
	}

	return &model.EncryptedReceipt{
		ID:                ids.UUID(),
		UserID:            userId,
		CreatedAt:         time.Now().Format(time.RFC3339),
		ReceiptName:       encryption.EncryptField(receiptInput.ReceiptName, aesKey, iv),
		RecipientPhone:    encryption.EncryptField(receiptInput.RecipientPhone, aesKey, iv),
		RecipientName:     encryption.EncryptField(receiptInput.RecipientName, aesKey, iv),
		RecipientEmail:    encryption.EncryptField(&receiptInput.RecipientEmail, aesKey, iv),
		RecipientAddress:  encryption.EncryptField(receiptInput.RecipientAddress, aesKey, iv),
		ReceiptNo:         receiptInput.ReceiptNo,
		Date:              receiptInput.Date,
		PaymentMethod:     stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(receiptInput.PaymentMethod), aesKey, iv)),
		PaymentNote:       encryption.EncryptField(receiptInput.PaymentNote, aesKey, iv),
		TotalAmount:       encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", totalAmount)), aesKey, iv),
		SubTotalAmount:    encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", subtotal)), aesKey, iv),
		TaxAmount:         encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", taxAmount)), aesKey, iv),
		EncryptedServices: make([]*model.EncryptedService, 0),
		AesIv:             stringUtil.StrPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:   stringUtil.StrPtr(string(encryption.EncryptKey(publicKeyPEM, aesKey))),
	}, nil
}

func whatsAppInputToReceiptModel(input model.SendReceiptPDFToWhatsApp, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	return &model.Receipt{
		ID:               ids.UUID(),
		UserID:           userId,
		CreatedAt:        time.Now().Format(time.RFC3339),
		ReceiptName:      receiptInput.ReceiptName,
		RecipientPhone:   &receiptInput.RecipientPhone,
		RecipientName:    receiptInput.RecipientName,
		RecipientEmail:   receiptInput.RecipientEmail,
		RecipientAddress: receiptInput.RecipientAddress,
		ReceiptNo:        receiptInput.ReceiptNo,
		Date:             receiptInput.Date,
		PaymentMethod:    receiptInput.PaymentMethod,
		PaymentNote:      receiptInput.PaymentNote,
		TotalAmount:      &totalAmount,
		SubTotalAmount:   &subtotal,
		TaxAmount:        &taxAmount,
		Services:         make([]*model.Service, 0),
	}
}

func whatsAppInputToEncryptedReceiptModel(input model.SendReceiptPDFToWhatsApp, userId string, totalAmount, subtotal, taxAmount float64, publicKeyPEM string) (*model.EncryptedReceipt, error) {
	receiptInput := input

	aesKey, iv, err := encryption.GenerateAESKeyAndIV()
	if err != nil {
		return nil, err
	}

	return &model.EncryptedReceipt{
		ID:                ids.UUID(),
		UserID:            userId,
		CreatedAt:         time.Now().Format(time.RFC3339),
		ReceiptName:       encryption.EncryptField(receiptInput.ReceiptName, aesKey, iv),
		RecipientPhone:    encryption.EncryptField(&receiptInput.RecipientPhone, aesKey, iv),
		RecipientName:     encryption.EncryptField(receiptInput.RecipientName, aesKey, iv),
		RecipientEmail:    encryption.EncryptField(receiptInput.RecipientEmail, aesKey, iv),
		RecipientAddress:  encryption.EncryptField(receiptInput.RecipientAddress, aesKey, iv),
		ReceiptNo:         receiptInput.ReceiptNo,
		Date:              receiptInput.Date,
		PaymentMethod:     stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(receiptInput.PaymentMethod), aesKey, iv)),
		PaymentNote:       encryption.EncryptField(receiptInput.PaymentNote, aesKey, iv),
		TotalAmount:       encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", totalAmount)), aesKey, iv),
		SubTotalAmount:    encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", subtotal)), aesKey, iv),
		TaxAmount:         encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", taxAmount)), aesKey, iv),
		EncryptedServices: make([]*model.EncryptedService, 0),
		AesIv:             stringUtil.StrPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:   stringUtil.StrPtr(string(encryption.EncryptKey(publicKeyPEM, aesKey))),
	}, nil
}

func downloadInputToReceiptModel(input model.DownloadPDF, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	return &model.Receipt{
		ID:               ids.UUID(),
		UserID:           userId,
		CreatedAt:        time.Now().Format(time.RFC3339),
		ReceiptName:      receiptInput.ReceiptName,
		RecipientPhone:   receiptInput.RecipientPhone,
		RecipientName:    receiptInput.RecipientName,
		RecipientEmail:   receiptInput.RecipientEmail,
		RecipientAddress: receiptInput.RecipientAddress,
		ReceiptNo:        receiptInput.ReceiptNo,
		Date:             receiptInput.Date,
		PaymentMethod:    receiptInput.PaymentMethod,
		PaymentNote:      receiptInput.PaymentNote,
		TotalAmount:      &totalAmount,
		SubTotalAmount:   &subtotal,
		TaxAmount:        &taxAmount,
		Services:         make([]*model.Service, 0),
	}
}

func downlaodInputToEncryptedReceiptModel(input model.DownloadPDF, userId string, totalAmount, subtotal, taxAmount float64, publicKeyPEM string) (*model.EncryptedReceipt, error) {
	receiptInput := input

	aesKey, iv, err := encryption.GenerateAESKeyAndIV()
	if err != nil {
		return nil, err
	}

	return &model.EncryptedReceipt{
		ID:                ids.UUID(),
		UserID:            userId,
		CreatedAt:         time.Now().Format(time.RFC3339),
		ReceiptName:       encryption.EncryptField(receiptInput.ReceiptName, aesKey, iv),
		RecipientPhone:    encryption.EncryptField(receiptInput.RecipientPhone, aesKey, iv),
		RecipientName:     encryption.EncryptField(receiptInput.RecipientName, aesKey, iv),
		RecipientEmail:    encryption.EncryptField(receiptInput.RecipientEmail, aesKey, iv),
		RecipientAddress:  encryption.EncryptField(receiptInput.RecipientAddress, aesKey, iv),
		ReceiptNo:         receiptInput.ReceiptNo,
		Date:              receiptInput.Date,
		PaymentMethod:     stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(receiptInput.PaymentMethod), aesKey, iv)),
		PaymentNote:       encryption.EncryptField(receiptInput.PaymentNote, aesKey, iv),
		TotalAmount:       encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", totalAmount)), aesKey, iv),
		SubTotalAmount:    encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", subtotal)), aesKey, iv),
		TaxAmount:         encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", taxAmount)), aesKey, iv),
		EncryptedServices: make([]*model.EncryptedService, 0),
		AesIv:             stringUtil.StrPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:   stringUtil.StrPtr(string(encryption.EncryptKey(publicKeyPEM, aesKey))),
	}, nil
}

func addServicesToReceipt(receiptModel *model.Receipt, services []*model.CreateBulkService) {
	// Create Services
	for _, serviceInput := range services {
		serviceModel := &model.Service{
			ID:          ids.UUID(),
			ReceiptID:   receiptModel.ID,
			CreatedAt:   time.Now().Format(time.RFC3339),
			Description: serviceInput.Description,
			Rate:        serviceInput.Rate,
			Quantity:    serviceInput.Quantity,
			Amount:      serviceInput.Amount,
		}
		receiptModel.Services = append(receiptModel.Services, serviceModel)
	}
}

func checkUserMode(mode string, useCount int) error {
	if mode == trial && useCount >= trialUseCountLimit {
		return errors.New("trial limit exceeded")
	}
	if mode == starter && useCount >= starterUseCountLimit {
		return errors.New("starter limit exceeded")
	}
	if mode == growth && useCount >= growthUseCountLimit {
		return errors.New("growth limit exceeded")
	}
	return nil
}

func ParseStringToFloat64Ptr(s *string) (*float64, error) {
	if s == nil {
		return nil, nil
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func parseRSAPrivateKeyFromPEM(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not an RSA private key")
	}

	return rsaKey, nil
}

func canonicalJSON(v any) ([]byte, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return json.Marshal(v)
	}
	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyBytes, _ := json.Marshal(k)
		buf.Write(keyBytes)
		buf.WriteByte(':')
		valBytes, _ := canonicalJSON(m[k])
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
