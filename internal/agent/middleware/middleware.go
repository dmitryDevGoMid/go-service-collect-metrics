package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/cryptohashsha"
	"github.com/go-resty/resty/v2"
)

type CleintInterface interface {
	OnBeforeRequest()
	OnAfterResponse()
	AssimEncryptBody()
}

type Client struct {
	client     *resty.Client
	cfg        *config.Config
	hashSha256 cryptohashsha.HashSha256
	encrypt    asimencrypt.AsimEncrypt
}

func NewClientMiddleware(client *resty.Client, cfg *config.Config, sha256 cryptohashsha.HashSha256, encrypt asimencrypt.AsimEncrypt) CleintInterface {
	return &Client{client: client, cfg: cfg, hashSha256: sha256, encrypt: encrypt}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

// Assimetric Encrypt Decode by Private Key
func (cl *Client) AssimEncryptBody() {
	cl.client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		if cl.cfg.PathEncrypt.KeyEncryptEnbled {

			body := fmt.Sprintf("%s", req.Body)
			bodyEncrypt, err := cl.encrypt.Encrypt(string(body))

			if err != nil {
				log.Println("Error decrypt into middleware", err)
			} else {
				req.Body = string(bodyEncrypt)
			}
		}
		return nil
	})
}

func (cl *Client) OnBeforeRequest() {
	// Registering Request Middleware
	cl.client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		body := fmt.Sprintf("%s", req.Body)

		hashString, err := cl.hashSha256.GetSha256ByData([]byte(body))

		if err == nil {
			fmt.Println("Тело запроса--------->>>>>>>>>>:", body)
			fmt.Println("Записали заголовок HashSHA256--------->>>>>>>>>>:", hashString)
			c.SetHeader("HashSHA256", fmt.Sprintf("%x", hashString))
		}

		//c.R().SetBody(body)

		c.SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetHeader("Accept-Encoding", "gzip")
		// Проверяем конфигурацию по умолчанию идет сжатие GZIP
		//if config.cfg.Gzip.Enable {
		//	c.SetHeader("Content-Encoding", "gzip")
		//}

		// Now you have access to Client and current Request object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})
}

func (cl *Client) OnAfterResponse() {
	// Registering Response Middleware
	cl.client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need
		fmt.Println("RESPONSE:", resp.Header().Get("Content-Type"))
		fmt.Println(string(resp.Body()))

		return nil // if its success otherwise return error
	})
}
