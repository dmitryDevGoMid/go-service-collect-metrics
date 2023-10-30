package routes

import (
	"bytes"
	"io"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Миделвари для заголовков
func WriteContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() != "/" {
			if c.Request.Header.Get("Content-Type") != "application/json" {
				c.Writer.Header().Set("Content-Type", "text/plain")
			} else {
				c.Writer.Header().Set("Content-Type", "application/json")
			}
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

// Логируем request and response
func LoggerMiddleware(appLogger *logger.APILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		strat := time.Now()

		uri := c.Request.RequestURI
		method := c.Request.Method
		content := c.Request.Header.Get("Content-Type")
		body, _ := io.ReadAll(c.Request.Body)

		bodyString := string(body)

		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		c.Next()

		duration := time.Since(strat)

		appLogger.Infof(
			"uri %s method %s duration %s status %d size %d content %s ",
			uri,
			method,
			duration,
			c.Writer.Status(),
			c.Writer.Size(),
			content,
			bodyString,
		)
	}
}
