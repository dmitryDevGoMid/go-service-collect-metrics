package handlers

import (
	"fmt"
	"io"
	"strconv"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/compress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/decompress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/serialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"

	"net/http"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/restutils"

	"github.com/gin-gonic/gin"
)

/*
	Unmarshal - строка в структуру
	Marshal - структруа в строку
*/

// Интерфейс для обработчиков запросов
type MetricsHandlers interface {

	//######### NOT JSON ###########
	GetMetricsGauge(c *gin.Context)
	GetMetricsCounter(c *gin.Context)
	UpdateGauge(c *gin.Context)
	UpdateCounter(c *gin.Context)
	//######### NOT JSON ###########

	GetMetrics(c *gin.Context)
	UpdateMetrics(c *gin.Context)

	GetAllMetricsHTML(c *gin.Context)

	Update(c *gin.Context)
	Value(c *gin.Context)

	UpdatePostJSON(c *gin.Context)
	ValuePostJSON(c *gin.Context)
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

// ####################### POST NOT JSON ######################
// endPointsMetricsHandlers GetMetricsGauge
func (h *metricsHandlers) GetMetricsGauge(c *gin.Context) {
	metricName := c.Param("metric")

	resp, err := h.metricsRepository.GetMetricGauge(metricName)

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// endPointsMetricsHandlers GetMetricsCounter
func (h *metricsHandlers) GetMetricsCounter(c *gin.Context) {
	metricName := c.Param("metric")

	resp, err := h.metricsRepository.GetMetricCounter(metricName)

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

// End Points MetricsHandlers UpdateGauge
func (h *metricsHandlers) UpdateGauge(c *gin.Context) {
	metricName := c.Param("metric")

	metricValue, err := strconv.ParseFloat(c.Param("value"), 64)
	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, `Неверный параметр метрики!`)
		return
	}

	h.metricsRepository.UpdateMetricGauge(metricName, metricValue)

	c.Status(http.StatusOK)
}

// End Points MetricsHandlers UpdateCounter
func (h *metricsHandlers) UpdateCounter(c *gin.Context) {

	metric := c.Param("metric")
	value := c.Param("value")

	fmt.Println("Получили:", value)

	metricValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, `Неверный параметр метрики!`)
		return
	}

	fmt.Println("Типизация:", metricValue)

	h.metricsRepository.UpdateMetricCounter(metric, metricValue)

	c.Status(http.StatusOK)
}

//####################### POST NOT JSON ######################

//######################### POST JSON ######################

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

	body, err := checkGzip(c)
	//body, err := io.ReadAll(c.Request.Body)

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
func (h *metricsHandlers) serializerResponse(metricsSData *serialize.Metrics) string {

	serializer := serialize.NewSerializer(h.cfg)

	var sendStringMetrics string

	serializeErr := serializer.SetData(metricsSData).GetData(&sendStringMetrics)

	if serializeErr.Errors() != nil {
		panic(serializeErr.Errors().Error())
	}

	return sendStringMetrics

}

// endPointsMetricsHandlers GetMetrics
func (h *metricsHandlers) GetMetrics(c *gin.Context) {
	metrics := h.unSerializerRequest(c)

	// В конце закрываем запрос
	//defer c.Request.Body.Close()

	fmt.Println(metrics)

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

	//c.Data(http.StatusOK, "application/json", []byte(sendData))
	gZipAccept([]byte(sendData), c)
}

// endPointsMetricsHandlers UpdateMetrics
func (h *metricsHandlers) UpdateMetrics(c *gin.Context) {
	metrics := h.unSerializerRequest(c)

	// В конце закрываем запрос
	//defer c.Request.Body.Close()

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

	//c.Data(http.StatusOK, "application/json", []byte(sendData))
	gZipAccept([]byte(sendData), c)
}

func (h *metricsHandlers) UpdatePostJSON(c *gin.Context) {
	h.UpdateMetrics(c)
}

// Point Update
func (h *metricsHandlers) Update(c *gin.Context) {
	typeMetric := c.Param("type")

	switch val := typeMetric; val {
	case "gauge":
		h.UpdateGauge(c)
	case "counter":
		h.UpdateCounter(c)
	default:
		c.Status(http.StatusBadRequest)
	}
}

func (h *metricsHandlers) ValuePostJSON(c *gin.Context) {
	h.GetMetrics(c)
}

// Point Value
func (h *metricsHandlers) Value(c *gin.Context) {

	fmt.Println("VALUE Content-Type NOT JSON")
	typeMetric := c.Param("type")

	switch val := typeMetric; val {
	case "gauge":
		h.GetMetricsGauge(c)
	case "counter":
		h.GetMetricsCounter(c)
	default:
		c.Status(http.StatusBadRequest)
	}
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

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func gZipAccept(data []byte, c *gin.Context) {
	compress_ := false

	content := c.Request.Header.Values("Accept-Encoding")

	for _, val := range content {
		if val == "gzip" {
			compress_ = true
		}
	}

	if compress_ {
		c.Writer.Header().Set("Content-Encoding", "gzip")
		dataCompress, _ := compress.CompressGzip([]byte(data))
		c.Data(http.StatusOK, "application/json", dataCompress)
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(data))
}
