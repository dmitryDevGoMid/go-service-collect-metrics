package routes

import (
	"github.com/gin-gonic/gin"
)

// Миделвари для заголовков
func WriteContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() != "/" {
			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.Next()
		} else {
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			c.Next()
		}
	}
}
