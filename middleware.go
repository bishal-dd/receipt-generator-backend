package main

import (
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	//   ctx := c.Request.Context()
	  
	//   // Log the authorization header to debug
	//   fmt.Print(c.GetHeader("authorization"), "\n")
	//   authHeader := c.GetHeader("authorization")

	//   if authHeader == "" {
	// 	fmt.Print("Authorization header not found\n")
	// 	c.AbortWithStatusJSON(401, gin.H{"error": "Missing authorization token"})
	// 	return
	//   }

	//   claims, ok := clerk.SessionClaimsFromContext(ctx)

	//   if !ok {
	// 	fmt.Print("Invalid or missing session claims\n")
	// 	c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or missing session claims"})
	// 	return
	//   }
	  
	//   usr, err := user.Get(ctx, claims.Subject)
	//   if err != nil {
	// 	c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("User retrieval failed: %v", err)})
	// 	return
	//   }
	  
	//   if usr == nil {
	// 	c.AbortWithStatusJSON(401, gin.H{"error": "User does not exist"})
	// 	return
	//   }

	userId := "user_2pOGOF48y7Y7wHWs33WK0AZIROH"
	  ctxs := contextUtil.SetContextValue(c.Request.Context(), contextUtil.UserIDKey, userId)
	  ctxs = contextUtil.SetContextValue(ctxs, contextUtil.GinContextKey, c)
	  c.Request = c.Request.WithContext(ctxs)
	  c.Next()
	}
  }