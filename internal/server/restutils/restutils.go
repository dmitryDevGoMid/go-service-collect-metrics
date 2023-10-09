package restutils

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message"`
}

// возвращаем ошибку со статусом ответа
func GinWriteError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, response{message})
}
