package main

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/helper/jwtUtil"
	"github.com/gin-gonic/gin"
)
type ContextKey string

const (
	GinContextKey ContextKey = "GinContextKey"
	UserIDKey     ContextKey = "UserID"
)
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := jwtUtil.TokenFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		token, err := jwtUtil.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		claims, err := jwtUtil.ClaimsFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// Extract the user_id from claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "user_id not found in token"})
			return
		}
		// Inject the userID and Gin context into the request context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, GinContextKey, c)

		// Update the request with the new context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
