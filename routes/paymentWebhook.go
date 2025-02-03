package routes

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookPayload struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			SubscriptionID int    `json:"subscription_id"`
			CustomerID     int    `json:"customer_id"`
			BillingReason  string `json:"billing_reason"`
			UserEmail      string `json:"user_email"`
			Status         string `json:"status"`
			TotalFormatted string `json:"total_formatted"`
		} `json:"attributes"`
	} `json:"data"`
}


func verifySignature(c *gin.Context, body []byte) bool {
	signingSecret := os.Getenv("LEMON_SQUEEZY_SIGNING_SECRET")
	if signingSecret == "" {
		return false
	}

	signature := c.GetHeader("X-Signature")
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write(body)
	computedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(computedSignature))
}

func PaymentWebhook(c *gin.Context, db *gorm.DB) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if !verifySignature(c, body) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	mode := determineMode(payload.Data.Attributes.TotalFormatted)
	if err := updateUserMode(payload.Data.Attributes.UserEmail, mode, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user mode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}

func determineMode(total_formatted string) string {
	switch total_formatted {
	case "$4.99":
		return "paid"
	case "$3.99":
		return "starter"
	case "$9.99":
		return "growth"
	case "$19.99":
		return "business"
	default:
		return "trial"
	}
}

func updateUserMode(email string, mode string, db *gorm.DB) error {

	ctx := context.Background()
	// Fetch user from Clerk using email
	params := &user.ListParams{
		EmailAddresses: []string{email},
	}
	users, err := user.List(ctx, params)
	if err != nil {
		log.Printf("Error fetching user from Clerk: %v", err)
		return errors.New("no user found")
	}
	if len(users.Users) == 0 {
		log.Printf("No Clerk user found with email: %s", email)
		return errors.New("no user found")
	}
	clerkUser := users.Users[0]
	user := &model.User{
		ID: clerkUser.ID,
	}
	
	if err := db.Model(user).Updates(map[string]interface{}{"mode": mode}).Error; err != nil {
		return  err
	}
	return nil
}
