package sandlers

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`            // Параметр кодирую строкой, принося производительность в угоду наглядности.
	Delta *int64   `json:"delta,omitempty"` // counter
	Value *float64 `json:"value,omitempty"` // gauge
}

func TestGaugeGzipHandlers(t *testing.T) {
	//errRedirectBlocked := errors.New("HTTP redirect blocked")
	/*redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})*/
	httpc := resty.New().SetBaseURL("http://localhost:8080/")
	//SetHostURL("http://localhost:8080")
	//SetRedirectPolicy(redirPolicy)

	id := "GetSetZip" + strconv.Itoa(rand.Intn(256))

	t.Run("update", func(t *testing.T) {
		value := rand.Float64() * 1e6
		req := httpc.R().
			SetHeader("Hash", "none").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Type", "application/json")

		resp, err := req.SetBody(
			&Metrics{
				ID:    id,
				MType: "gauge",
				Value: &value,
			}).Post("update")

		dumpErr := assert.NoError(t, err,
			"Ошибка при попытке сделать запрос с обновлением gauge")
		/*dumpErr = dumpErr && assert.Equalf(t, resp.StatusCode(),
		"Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)*/
		fmt.Println(dumpErr)

		var result Metrics
		resp, err = req.
			SetBody(&Metrics{
				ID:    id,
				MType: "gauge",
			}).
			SetResult(&result).
			Post("value/")
		fmt.Println(resp)
		fmt.Println(err)
		//assert.Equalf()

		//assert.Equalf( resp.StatusCode(), resp.StatusCode())
		/*dumpErr = dumpErr && assert.NoError(t, err,
			"Ошибка при попытке сделать запрос с получением значения gauge")
		dumpErr = dumpErr && assert.Equalf(t, resp.StatusCode(), resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
		dumpErr = dumpErr && assert.Containsf(resp.Header().Get("Content-Type"), "application/json",
			"Заголовок ответа Content-Type содержит несоответствующее значение")
		dumpErr = dumpErr && assert.Containsf(t, "gzip","Заголовок ответа Content-Encoding содержит несоответствующее значение")
		dumpErr = dumpErr && assert.NotEqualf(nil, result.Value,
			"Несоответствие отправленного значения gauge (%f) полученному от сервера (nil), '%q %s'", value, req.Method, req.URL)
		dumpErr = dumpErr && assert.Equalf(t, resp.StatusCode(), *result.Value,
			"Несоответствие отправленного значения gauge (%f) полученному от сервера (%f), '%q %s'", value, *result.Value, req.Method, req.URL)
		// dumpErr = dumpErr && suite.Equal(suite.Hash(&result), result.Hash, "Хеш-сумма не соответствует расчетной")

		/*if !dumpErr {
			dump := dumpRequest(req.RawRequest, true)
			suite.T().Logf("Оригинальный запрос:\n\n%s", dump)
			dump = dumpResponse(resp.RawResponse, true)
			suite.T().Logf("Оригинальный ответ:\n\n%s", dump)
		}*/
	})
}
