package main

import (
	"fmt"
	"strings"

	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
        
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

        // Set user ID in context
        ctxWithUser := contextUtil.SetContextValue(ctx, contextUtil.UserIDKey, usr.ID)
        ctxWithGin := contextUtil.SetContextValue(ctxWithUser, contextUtil.GinContextKey, c)
        c.Request = c.Request.WithContext(ctxWithGin)

        c.Next()
	}
  }