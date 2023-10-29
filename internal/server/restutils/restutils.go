package restutils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyBody    = errors.New("empty body")
	ErrUnauthorized = errors.New("unauthorized")
)

type response struct {
	Message string `json:"message"`
}

// возвращаем ошибку со статусом ответа
func GinWriteError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}
