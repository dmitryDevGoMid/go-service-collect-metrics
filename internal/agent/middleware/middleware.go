package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/go-resty/resty/v2"
)

type CleintInterface interface {
	OnBeforeRequest()
	OnAfterResponse()
}

type Client struct {
	client *resty.Client
	cfg    *config.Config
}

func NewClientMiddleware(client *resty.Client, cfg *config.Config) CleintInterface {
	return &Client{client: client, cfg: cfg}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func (config *Client) OnBeforeRequest() {
	// Registering Request Middleware
	config.client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		c.SetHeader("Content-Type", "application/json")

		// Проверяем конфигурацию по умолчанию идет сжатие GZIP
		if config.cfg.Gzip.Enable {
			//c.SetHeader("Content-Encoding", "gzip")
		}

		// Now you have access to Client and current Request object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})
}

func (config *Client) OnAfterResponse() {
	// Registering Response Middleware
	config.client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need
		fmt.Println("RESPONSE:", resp.Header().Get("Content-Type"))
		fmt.Println(string(resp.Body()))

		return nil // if its success otherwise return error
	})
}
