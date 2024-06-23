package loaders

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Middleware injects data loaders into the context
func LoaderMiddleware(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		loader := NewLoaders(conn)
		ctx := context.WithValue(c.Request.Context(), loadersKey, loader)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}