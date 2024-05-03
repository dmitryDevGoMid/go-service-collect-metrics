// Package handlers provides handlers for handling metrics.
//
// Metrics are stored in a repository and can be retrieved and updated using the handlers.
// The handlers support both JSON and non-JSON formats.
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
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mrepository"

	"net/http"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/restutils"

	"github.com/gin-gonic/gin"
)

/*
	Unmarshal - строка в структуру
	Marshal - структруа в строку
*/

// MetricsHandlers is an interface for metrics handlers.
type MetricsHandlers interface {
	// GetMetricsGauge retrieves the gauge metric value for the given metric name.
	// It returns the value as a string in plain text format.
	GetMetricsGauge(c *gin.Context)

	// GetMetricsCounter retrieves the counter metric value for the given metric name.
	// It returns the value as a string in plain text format.
	GetMetricsCounter(c *gin.Context)

	// UpdateGauge updates the gauge metric value for the given metric name.
	// It accepts the value as a float64 parameter in the URL.
	UpdateGauge(c *gin.Context)

	// UpdateCounter updates the counter metric value for the given metric name.
	// It accepts the value as an int64 parameter in the URL.
	UpdateCounter(c *gin.Context)

	// Updates updates multiple metrics at once using a POST request with JSON data.
	Updates(c *gin.Context)

	// GetMetrics retrieves all metrics in JSON format.
	GetMetrics(c *gin.Context)

	// UpdateMetrics updates multiple metrics at once using a PUT request with JSON data.
	UpdateMetrics(c *gin.Context)

	// GetAllMetricsHTML retrieves all metrics in HTML format.
	GetAllMetricsHTML(c *gin.Context)

	// Update updates a single metric using a POST request with form data.
	Update(c *gin.Context)

	// Value retrieves the value of a single metric using a GET request with form data.
	Value(c *gin.Context)

	// UpdatePostJSON updates a single metric using a POST request with JSON data.
	UpdatePostJSON(c *gin.Context)

	// ValuePostJSON retrieves the value of a single metric using a GET request with JSON data.
	ValuePostJSON(c *gin.Context)

	// Ping checks the availability of the service.
	Ping(c *gin.Context)
}

// metricsHandlers is a struct that implements the MetricsHandlers interface.
type MetricsHandlersType struct {
	managerRepository mrepository.ManagerRepository
	metricsRepository repository.MetricsRepository
	cfg               *config.Config
}

// Конструктор
/*func NewMetricsHandlers(
	managerRepository mrepository.ManagerRepository,
	metricsRepository repository.MetricsRepository,
	cfg *config.Config) MetricsHandlers {
	return &metricsHandlers{metricsRepository: metricsRepository, managerRepository: managerRepository, cfg: cfg}
}*/

// NewMetricsHandlers creates a new instance of metricsHandlers.
func NewMetricsHandlers(
	managerRepository mrepository.ManagerRepository,
	metricsRepository repository.MetricsRepository,
	cfg *config.Config) MetricsHandlers {
	return &MetricsHandlersType{metricsRepository: metricsRepository, managerRepository: managerRepository, cfg: cfg}
}

// Change repository
func (h *MetricsHandlersType) setRepository() {
	//Если менеджер репозитария пустой, то используем репозитарий назначенный через конструктор
	if h.managerRepository != nil {
		h.metricsRepository = h.managerRepository.GetRepositoryActive()
	}
}

// ####################### POST NOT JSON ######################
// End Points MetricsHandlers GetMetricsGauge
// Returns the current value of a gauge metric by its name.
// Parameters:
// - metric: the name of the metric
// Returns:
// - HTTP status code 200 OK and the metric value as a string in the response body if the metric is found
// - HTTP status code 404 Not Found if the metric is not found
func (h *MetricsHandlersType) GetMetricsGauge(c *gin.Context) {
	h.setRepository()

	metricName := c.Param("metric")

	retryTest := repository.Decorator{IMetric: h.metricsRepository}
	resp, err := retryTest.GetMetricGauge(c, metricName)

	respString := fmt.Sprintf("%v", resp)

	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.Data(http.StatusOK, "text/plain", []byte(respString))
	}
}

// End Points MetricsHandlers GetMetricsCounter
// Returns the value of a counter metric by its name.
// Parameters:
// - metric: the name of the metric
// Returns:
// - HTTP status code 200 OK and the metric value as a string in the response body if the metric is found
// - HTTP status code 404 Not Found if the metric is not found
func (h *MetricsHandlersType) GetMetricsCounter(c *gin.Context) {
	h.setRepository()

	metricName := c.Param("metric")

	retryTest := repository.Decorator{IMetric: h.metricsRepository}
	resp, err := retryTest.GetMetricCounter(c, metricName)

	respString := fmt.Sprintf("%d", resp)

	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.Data(http.StatusOK, "text/plain", []byte(respString))
	}
}

// End Points MetricsHandlers UpdateGauge
// Updates the value of a gauge metric by its name.
// Parameters:
// - metric: the name of the metric
// - value: the new value of the metric as a floating-point number
// Returns:
// - HTTP status code 200 OK on successful execution
// - HTTP status code 400 Bad Request if the metric value is invalid
func (h *MetricsHandlersType) UpdateGauge(c *gin.Context) {
	h.setRepository()

	metricName := c.Param("metric")

	metricValue, err := strconv.ParseFloat(c.Param("value"), 64)
	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, `Неверный параметр метрики!`)
		return
	}

	h.metricsRepository.UpdateMetricGauge(c, metricName, metricValue)

	c.Status(http.StatusOK)
}

// End Points MetricsHandlers UpdateCounter
// Updates the value of a counter metric by its name.
// Parameters:
// - metric: the name of the metric
// - value: the new value of the metric as an integer
// Returns:
// - HTTP status code 200 OK on successful execution
// - HTTP status code 400 Bad Request if the metric value is invalid
func (h *MetricsHandlersType) UpdateCounter(c *gin.Context) {
	h.setRepository()
	metric := c.Param("metric")
	value := c.Param("value")

	fmt.Println("Получили:", value)

	metricValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, `Неверный параметр метрики!`)
		return
	}

	fmt.Println("Типизация:", metricValue)

	h.metricsRepository.UpdateMetricCounter(c, metric, metricValue)

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
func (h *MetricsHandlersType) unSerializerRequestBatch(c *gin.Context) []unserialize.Metrics {
	if c.Request.Body == nil {
		restutils.GinWriteError(c, http.StatusBadRequest, restutils.ErrEmptyBody.Error())
		return []unserialize.Metrics{}
	}

	body, err := checkGzip(c)
	//body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, err.Error())
		return []unserialize.Metrics{}
	}

	var metrics []unserialize.Metrics

	unserializeData := unserialize.NewUnSerializer(h.cfg)

	unserializeError := unserializeData.SetData(&body).GetDataBatch(&metrics)

	if unserializeError.Errors() != nil {
		panic(unserializeError.Errors().Error())
	}

	return metrics
}

// Point Serialize Data by Request
func (h *MetricsHandlersType) unSerializerRequest(c *gin.Context) unserialize.Metrics {
	if c.Request.Body == nil {
		restutils.GinWriteError(c, http.StatusBadRequest, restutils.ErrEmptyBody.Error())
		return unserialize.Metrics{}
	}

	body, err := checkGzip(c)
	//body, err := io.ReadAll(c.Request.Body)

	//fmt.Println(string(body))

	if err != nil {
		restutils.GinWriteError(c, http.StatusBadRequest, err.Error())
		return unserialize.Metrics{}
	}

	var metrics unserialize.Metrics

	unserializeData := unserialize.NewUnSerializer(h.cfg)

	unserializeError := unserializeData.SetData(&body).GetData(&metrics)

	if unserializeError.Errors() != nil {
		//panic(unserializeError.Errors().Error())
		fmt.Println("UnserializeError=>", unserializeError.Errors().Error())
	}

	//fmt.Println(metrics)

	return metrics
}

// Point Serialize Data for Send
func (h *MetricsHandlersType) serializerResponse(metricsSData *serialize.Metrics) string {

	serializer := serialize.NewSerializer(h.cfg)

	var sendStringMetrics string

	serializeErr := serializer.SetData(metricsSData).GetData(&sendStringMetrics)

	if serializeErr.Errors() != nil {
		panic(serializeErr.Errors().Error())
	}

	return sendStringMetrics

}

// GetMetrics retrieves the value of a metric by its ID and type.
// It accepts a request body containing the metric ID and type, and returns the metric value in the response body.
// The type can be either "gauge" or "counter". If the type is not valid, it returns a 400 Bad Request status code.
// If the metric is not found, it returns a 404 Not Found status code.
// The response body is compressed using gzip if the client accepts it.
func (h *MetricsHandlersType) GetMetrics(c *gin.Context) {
	h.setRepository()

	metrics := h.unSerializerRequest(c)

	// В конце закрываем запрос
	//defer c.Request.Body.Close()

	//fmt.Println(metrics)

	if metrics == (unserialize.Metrics{}) {
		return
	}

	typeMetric := metrics.MType

	var respCounter int64
	var respGauge float64
	var err error

	switch val := typeMetric; val {
	case "gauge":
		retry := repository.Decorator{IMetric: h.metricsRepository}
		respGauge, err = retry.GetMetricGauge(c, metrics.ID)
		//respGauge, err = h.metricsRepository.GetMetricGauge(c, metrics.ID)
	case "counter":
		respCounter, err = h.metricsRepository.GetMetricCounter(c, metrics.ID)
		//retry := repository.Decorator{IMetric: h.metricsRepository}
		//respCounter, err = retry.GetMetricCounter(c, metrics.ID)
	default:
		c.Status(http.StatusBadRequest)
		return
	}

	if err != nil {
		c.Status(http.StatusNotFound)
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
	data := gZipAccept([]byte(sendData), c)
	c.Data(http.StatusOK, "application/json", []byte(data))
	return
}

// UpdateMetrics updates the value of a metric by its ID and type.
// It accepts a request body containing the metric ID, type, and new value or delta.
// The type can be either "gauge" or "counter". If the type is not valid, it returns a 400 Bad Request status code.
// If the metric is not found, it returns a 404 Not Found status code.
// The updated metric value is returned in the response body.
// The response body is compressed using gzip if the client accepts it.
func (h *MetricsHandlersType) UpdateMetrics(c *gin.Context) {
	h.setRepository()

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
		h.metricsRepository.UpdateMetricGauge(c, metrics.ID, *metrics.Value)
		respGauge, err = h.metricsRepository.GetMetricGauge(c, metrics.ID)
	case "counter":
		h.metricsRepository.UpdateMetricCounter(c, metrics.ID, *metrics.Delta)
		respCounter, err = h.metricsRepository.GetMetricCounter(c, metrics.ID)
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
	data := gZipAccept([]byte(sendData), c)
	c.Data(http.StatusOK, "application/json", []byte(data))
}

func (h *MetricsHandlersType) UpdatePostJSON(c *gin.Context) {
	h.UpdateMetrics(c)
}

// Point Update
func (h *MetricsHandlersType) Update(c *gin.Context) {
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

func (h *MetricsHandlersType) ValuePostJSON(c *gin.Context) {
	h.GetMetrics(c)
}

// Value retrieves the value of a metric by its type and name.
// It accepts a path parameter "type" which can be either "gauge" or "counter" to specify the type of the metric.
// If the type is not valid, it returns a 400 Bad Request status code.
// Otherwise, it calls the corresponding handler function to retrieve the value of the metric and returns it in the response body.
func (h *MetricsHandlersType) Value(c *gin.Context) {

	//fmt.Println("VALUE Content-Type NOT JSON")
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
func (h *MetricsHandlersType) GetAllMetricsHTML(c *gin.Context) {
	h.setRepository()

	html := ""
	metrics, err := h.metricsRepository.GetAllMetrics(c)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

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
		dataCompress, err := compress.CompressGzip([]byte(html))
		if err != nil {
			fmt.Println("Error:", err)
		}
		c.Data(http.StatusOK, "", dataCompress)
		return
	}

	c.Data(http.StatusOK, "text/html", []byte(html))
}

func gZipAccept(data []byte, c *gin.Context) []byte {
	compress_ := false

	content := c.Request.Header.Values("Accept-Encoding")

	for _, val := range content {
		if val == "gzip" {
			compress_ = true
		}
	}

	if compress_ {
		c.Writer.Header().Set("Content-Encoding", "gzip")
		dataCompress, err := compress.CompressGzip(data)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return dataCompress
	}

	return data
}

// Ping Data Base Postgres Server
// The Ping function handles the HTTP request to check the connection to the Postgres server.
// It returns a 200 status code if the connection is successful, and a 500 status code if it fails.
func (h *MetricsHandlersType) Ping(c *gin.Context) {
	h.setRepository()

	err := h.metricsRepository.PingDatabase(c)

	if err != nil {
		c.Data(500, "Ping failed", []byte("Failed to ping database"))
		return
	}

	c.Data(200, "Ping successful", []byte("Success to ping database"))
}

// Updates
// The Updates function handles the HTTP request to update multiple metrics in the database.
// It expects a request body containing a list of metrics in JSON format, and saves them to the database.
// It returns a 200 status code if the update is successful, and a 400 status code if the request body is empty or invalid.
func (h *MetricsHandlersType) Updates(c *gin.Context) {
	metrics := h.unSerializerRequestBatch(c)

	err := h.metricsRepository.SaveMetricsBatch(c, metrics)

	if err != nil {
		fmt.Println("Error Save Metrics: ", err)
		restutils.GinWriteError(c, http.StatusBadRequest, restutils.ErrEmptyBody.Error())
		return
	}

	c.Data(200, "Updates successful", []byte("Success get to Updates"))
}
