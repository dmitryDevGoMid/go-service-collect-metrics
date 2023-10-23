package routes

import (
	"github.com/gin-gonic/gin"
)

// Миделвари для заголовков
func WriteContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() != "/" {
			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		} else {
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length")
		c.Header("Access-Control-Allow-Methods", "POST, GET")

		c.Next()
	}
}
