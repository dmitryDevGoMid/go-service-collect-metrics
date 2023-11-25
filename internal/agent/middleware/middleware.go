package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/cryptohashsha"
	"github.com/go-resty/resty/v2"
)

type CleintInterface interface {
	OnBeforeRequest()
	OnAfterResponse()
}

type Client struct {
	client     *resty.Client
	cfg        *config.Config
	hashSha256 cryptohashsha.HashSha256
}

func NewClientMiddleware(client *resty.Client, cfg *config.Config, sha256 cryptohashsha.HashSha256) CleintInterface {
	return &Client{client: client, cfg: cfg, hashSha256: sha256}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func (cl *Client) OnBeforeRequest() {
	// Registering Request Middleware
	cl.client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		body := fmt.Sprintf("%s", req.Body)
		//fmt.Println(body)

		hashString, err := cl.hashSha256.GetSha256ByData([]byte(body))

		if err == nil {
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
