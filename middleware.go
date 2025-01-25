package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var userId string // Declare `userId` as a variable for reassignment

		if os.Getenv("ENV") == "development" {
			userId = "user_2rq8tXKZAxaKJEcSti29HmQHmKt"
		} else {
			// Extract the Bearer token
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Missing authorization header"})
				return
			}

			// Remove "Bearer " prefix
			sessionToken := strings.TrimPrefix(authHeader, "Bearer ")
			if sessionToken == authHeader {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
				return
			}

			// Verify the JWT token
			claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
				Token: sessionToken,
			})
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
				return
			}

			// Get user details
			usr, err := user.Get(ctx, claims.Subject)
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("Failed to get user: %v", err)})
				return
			}
			userId = usr.ID
		}

		// Set user ID in context
		ctxWithUser := contextUtil.SetContextValue(ctx, contextUtil.UserIDKey, userId)
		ctxWithGin := contextUtil.SetContextValue(ctxWithUser, contextUtil.GinContextKey, c)
		c.Request = c.Request.WithContext(ctxWithGin)

		c.Next()
	}
}
