package handlers

import (
	"fmt"
	"strconv"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	//"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/service"
	"net/http"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/restutils"

	"github.com/gin-gonic/gin"
)

// Интерфейс для обработчиков запросов
type MetricsHandlers interface {
	GetMetricsGauge(c *gin.Context)
	GetMetricsCounter(c *gin.Context)
	UpdateGauge(c *gin.Context)
	UpdateCounter(c *gin.Context)
	GetAllMetricsHTML(c *gin.Context)
	Update(c *gin.Context)
	Value(c *gin.Context)
}

// Структура реализующая интерфейс
type metricsHandlers struct {
	metricsRepository repository.MetricsRepository
	//metricsService service.MetricsService
}

// Конструктор
func NewMetricsHandlers(metricsRepository repository.MetricsRepository) MetricsHandlers {
	return &metricsHandlers{metricsRepository: metricsRepository}
}

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

func (h *metricsHandlers) Value(c *gin.Context) {

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
	c.Data(http.StatusOK, "", []byte(html))
}
