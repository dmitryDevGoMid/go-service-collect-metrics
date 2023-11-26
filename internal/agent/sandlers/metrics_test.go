package sandlers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`            // Параметр кодирую строкой, принося производительность в угоду наглядности.
	Delta *int64   `json:"delta,omitempty"` //counter
	Value *float64 `json:"value,omitempty"` //gauge
	Hash  string   `json:"hash,omitempty"`  //counter
}

func TestUpdateGzipHandlers(t *testing.T) {

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	baseUrl := fmt.Sprintf("http://%s", cfg.Server.Address)

	httpc := resty.New().
		SetBaseURL(baseUrl).
		SetRedirectPolicy(redirPolicy)

	id := "GetGoodSetZipus" + strconv.Itoa(rand.Intn(256))

	t.Run("update", func(t *testing.T) {
		value := rand.Float64() * 1e6
		req := httpc.R().
			SetHeader("Hash", "none").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Type", "application/json")

		resp, err := httpc.R().SetBody(
			&Metrics{
				ID:    id,
				MType: "gauge",
				Value: &value,
			}).Post("update/")

		if err != nil {
			fmt.Println(err)
		}

		assert.Equal(t, resp.StatusCode(), 200)

		var result Metrics
		resp, err = req.
			SetBody(&Metrics{
				ID:    id,
				MType: "gauge",
			}).
			SetResult(&result).
			Post("value/")

		if err != nil {
			fmt.Println(err)
		}

		assert.Equal(t, resp.StatusCode(), 200)
		assert.Equal(t, *result.Value, value)
	})
}
