package routes

import (
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileRequest struct {
    UserID string `json:"userId"`
}

func AddProfile(c *gin.Context, db *gorm.DB) {
    var req ProfileRequest
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request body"})
        return
    }

    if req.UserID == "" {
        c.JSON(400, gin.H{"error": "Missing userId"})
        return
    }

	me := &model.User{
		ID: req.UserID,
		Mode: "trial",
		UseCount: 0,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := db.Create(me).Error; err != nil {
        c.JSON(400, gin.H{"error": "Cannot create user"})
		return
	}

    profile := model.Profile{
		ID: ids.UUID(),
        UserID: req.UserID,
		PhoneNumberCountryCode: "US",
		Currency: "USD",
		CreatedAt: time.Now().Format(time.RFC3339),
    }

    if err := db.Create(&profile).Error; err != nil {
        c.JSON(400, gin.H{"error": "Cannot create profile"})
        return
    }

    c.JSON(200, gin.H{"success": true, "profile": profile})
}