package routes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/compress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/cryptohashsha"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/decompress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/logger"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/file"
	"github.com/gin-gonic/gin"
)

// Миделвари для заголовков
func WriteContentType() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() != "/" {
			if c.Request.Header.Get("Accept") != "application/json" {
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

func checkHeader(c *gin.Context, nameHeader string, keyHeader string) bool {
	content := c.Request.Header.Values(nameHeader)

	isKey := false

	for _, val := range content {
		if val == keyHeader {
			isKey = true
		}
	}

	return isKey
}

// Передаем данные сжатые, если есть соответствующий заголовок
func ToolsGroupPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("dataCompress")
		if checkHeader(c, "Accept-Encoding", "gzip") {
			wb := &toolBodyWriter{
				body:           &bytes.Buffer{},
				ResponseWriter: c.Writer,
			}
			c.Writer = wb

			c.Next()

			c.Writer.Header().Set("Content-Encoding", "gzip")

			originBytes := wb.body
			fmt.Printf("%s", originBytes)
			fmt.Println("")

			// clear Origin Buffer
			wb.body = &bytes.Buffer{}
			//bodyString := obj.String("%s", originBytes.String())

			dataCompress, _ := gzipCompress(originBytes.Bytes())

			wb.Write(dataCompress)
			wb.ResponseWriter.Write(wb.body.Bytes())

			fmt.Println(string(dataCompress))
			fmt.Println("dataCompress")

			//fmt.Println("Replace BODY:", wb.body)

			//c.Data(http.StatusOK, "application/json", dataCompress)
		} else {
			c.Next()
		}
	}
}

// Интерфейс gin.ResponseWriter позволяет добраться до тела ответа и перезаписать его
type toolBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r toolBodyWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// Выполняем разархивирование прилитевших данных, если есть соответствующий заголовок
func DecompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if checkHeader(c, "Content-Encoding", "gzip") {

			body, _ := io.ReadAll(c.Request.Body)

			decompressBody, _ := checkGzipAndDecompress(body)

			c.Request.Body = io.NopCloser(bytes.NewReader(decompressBody))
		}

		c.Next()

	}
}

// Compress
func gzipCompress(body []byte) ([]byte, error) {

	return compress.CompressGzip(body)
}

// Decompress
func checkGzipAndDecompress(body []byte) ([]byte, error) {

	decompr, _ := decompress.DecompressGzip(body)

	fmt.Println("MIDDLE DECOMPRESS====>", decompr)
	fmt.Println("MIDDLE STRING DECOMPRESS====>", string(decompr))
	return decompr, nil
}

// Логируем request and response
func LoggerMiddleware(appLogger *logger.APILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		strat := time.Now()

		uri := c.Request.RequestURI
		method := c.Request.Method
		content := c.Request.Header.Get("Content-Type")
		//body := c.Request.Body
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		myString := string(body[:])

		c.Next()

		duration := time.Since(strat)

		appLogger.Infof(
			"uri %s method %s duration %s status %d size %d body %s",
			uri,
			method,
			duration,
			c.Writer.Status(),
			c.Writer.Size(),
			content,
			myString,
		)
	}
}

// Assimetric Encrypt Decode by Private Key
func AssimEncryptBody(cfg *config.Config, encrypt asimencrypt.AsimEncrypt) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.PathEncrypt.KeyEncryptEnbled {
			body, _ := io.ReadAll(c.Request.Body)
			fmt.Println("Получили--------->>>>>>>>>>:", body)
			bodyDecrypt, err := encrypt.Decrypt(body)
			fmt.Println("Расшифровали--------->>>>>>>>>>:", bodyDecrypt)

			if err != nil {
				log.Println("Error decrypt into middleware", err)
				bodyDecrypt = string(body)
			}

			c.Request.Body = io.NopCloser(bytes.NewReader([]byte(bodyDecrypt)))
		}
		c.Next()
	}
}

// Пишем синхронно с отправкой данные на диск в случае update данных
func CheckHashSHA256Data(cfg *config.Config, hash cryptohashsha.HashSha256) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println("===CheckHashSHA256Data===", cfg.HashSHA256.Key)

		hashData := c.Request.Header.Get("HashSHA256")
		fmt.Println("В заголовке есть ключ", hashData)

		if hashData != "" && hash.CheckHashSHA256Key() {
			body, _ := io.ReadAll(c.Request.Body)
			fmt.Println("Тело запроса======>", body)

			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			if !hash.CheckHashSHA256Data(body, hashData) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"message": "hash does not match",
				})
			}
		}

		c.Next()

	}
}

// Пишем синхронно с отправкой данные на диск в случае update данных
func SaveFileToDisk(config *config.Config, file file.WorkerFile) gin.HandlerFunc {
	return func(c *gin.Context) {

		beforePath := c.FullPath()
		afterPath := strings.Split(beforePath, "update")

		c.Next()

		if afterPath[0] == "/" {
			if c.Writer.Status() == 200 && config.File.StoreInterval == 0 {
				file.SaveAllMetrics(c)
			}
		}
	}
}
