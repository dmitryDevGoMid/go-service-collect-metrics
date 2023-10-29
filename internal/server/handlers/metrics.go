package handlers

import (
	"fmt"
	"io"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/compress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/decompress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/serialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/unserialize"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"

	"net/http"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/restutils"

	"github.com/gin-gonic/gin"
)

// Интерфейс для обработчиков запросов
type MetricsHandlers interface {
	GetMetrics(c *gin.Context)
	UpdateMetrics(c *gin.Context)

	GetAllMetricsHTML(c *gin.Context)

	Update(c *gin.Context)
	Value(c *gin.Context)
}

// Структура реализующая интерфейс
type metricsHandlers struct {
	metricsRepository repository.MetricsRepository
	cfg               *config.Config
}

// Конструктор
func NewMetricsHandlers(metricsRepository repository.MetricsRepository, cfg *config.Config) MetricsHandlers {
	return &metricsHandlers{metricsRepository: metricsRepository, cfg: cfg}
}

func checkGzip(c *gin.Context) ([]byte, error) {
	compress_ := false

	content := c.Request.Header.Values("Content-Encoding")

	for _, val := range content {
		if val == "gzip" {
			compress_ = true
		}
	}

	body, _ := io.ReadAll(c.Request.Body)

	if compress_ {

		decompr, _ := decompress.DecompressGzip(body)

		return decompr, nil
	}

	return body, nil
}

// Point Serialize Data by Request
func (h *metricsHandlers) unSerializerRequest(c *gin.Context) unserialize.Metrics {
	if c.Request.Body == nil {
		restutils.GinWriteError(c, http.StatusBadRequest, restutils.ErrEmptyBody.Error())
		return unserialize.Metrics{}
	}

	// В конце закрываем запрос
	defer c.Request.Body.Close()

	body, err := checkGzip(c)

	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, err.Error())
		return unserialize.Metrics{}
	}

	var metrics unserialize.Metrics

	unserializeData := unserialize.NewUnSerializer(h.cfg)

	unserializeError := unserializeData.SetData(&body).GetData(&metrics)

	if unserializeError.Errors() != nil {
		panic(unserializeError.Errors().Error())
	}

	return metrics
}

// Point Serialize Data for Send
func (h *metricsHandlers) serializerResponse(metricsSData *serialize.Metrics) *serialize.Metrics {

	serializer := serialize.NewSerializer(h.cfg)

	var sendStringMetrics string

	serialize_err := serializer.SetData(metricsSData).GetData(&sendStringMetrics)

	if serialize_err.Errors() != nil {
		panic(serialize_err.Errors().Error())
	}

	return metricsSData

}

// endPointsMetricsHandlers GetMetrics
func (h *metricsHandlers) GetMetrics(c *gin.Context) {
	metrics := h.unSerializerRequest(c)

	if metrics == (unserialize.Metrics{}) {
		return
	}

	typeMetric := metrics.MType

	var respCounter int64
	var respGauge float64
	var err error

	switch val := typeMetric; val {
	case "gauge":
		respGauge, err = h.metricsRepository.GetMetricGauge(metrics.ID)
	case "counter":
		respCounter, err = h.metricsRepository.GetMetricCounter(metrics.ID)
	default:
		c.Status(http.StatusBadRequest)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	var metricsSData serialize.Metrics

	switch val := typeMetric; val {
	case "gauge":
		metricsSData = serialize.Metrics{ID: metrics.ID, MType: typeMetric, Delta: nil, Value: &respGauge}
	case "counter":
		metricsSData = serialize.Metrics{ID: metrics.ID, MType: typeMetric, Delta: &respCounter, Value: nil}
	default:
		c.Status(http.StatusBadRequest)
	}

	sendData := h.serializerResponse(&metricsSData)

	c.JSON(http.StatusOK, sendData)
}

// endPointsMetricsHandlers UpdateMetrics
func (h *metricsHandlers) UpdateMetrics(c *gin.Context) {
	metrics := h.unSerializerRequest(c)

	if metrics == (unserialize.Metrics{}) {
		return
	}

	typeMetric := metrics.MType

	var respCounter int64
	var respGauge float64
	var err error

	switch val := typeMetric; val {
	case "gauge":
		h.metricsRepository.UpdateMetricGauge(metrics.ID, *metrics.Value)
		respGauge, err = h.metricsRepository.GetMetricGauge(metrics.ID)
	case "counter":
		h.metricsRepository.UpdateMetricCounter(metrics.ID, *metrics.Delta)
		respCounter, err = h.metricsRepository.GetMetricCounter(metrics.ID)
	default:
		c.Status(http.StatusBadRequest)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	var metricsSData serialize.Metrics

	switch val := typeMetric; val {
	case "gauge":
		metricsSData = serialize.Metrics{ID: metrics.ID, MType: typeMetric, Delta: nil, Value: &respGauge}
	case "counter":
		metricsSData = serialize.Metrics{ID: metrics.ID, MType: typeMetric, Delta: &respCounter, Value: nil}
	default:
		c.Status(http.StatusBadRequest)
	}

	sendData := h.serializerResponse(&metricsSData)

	c.JSON(http.StatusOK, sendData)
}

// Point Update
func (h *metricsHandlers) Update(c *gin.Context) {
	h.UpdateMetrics(c)
}

// Point Value
func (h *metricsHandlers) Value(c *gin.Context) {
	h.GetMetrics(c)
}

// End Points MetricsHandlers GetAllMetricsHtml
func (h *metricsHandlers) GetAllMetricsHTML(c *gin.Context) {
	html := ""
	metrics := h.metricsRepository.GetAllMetrics()
	for key, val := range metrics.Counter {
		html += fmt.Sprintf("<div>%s => %d </div>", key, val)
	}
	for key, val := range metrics.Gauge {
		html += fmt.Sprintf("<div>%s => %v </div>", key, val)
	}

	compress_ := false

	content := c.Request.Header.Values("Accept-Encoding")

	for _, val := range content {
		if val == "gzip" {
			compress_ = true
		}
	}

	if compress_ {
		c.Writer.Header().Set("Content-Encoding", "gzip")
		dataCompress, _ := compress.CompressGzip([]byte(html))
		c.Data(http.StatusOK, "", dataCompress)
		return
	}

	c.Data(http.StatusOK, "", []byte(html))
}
