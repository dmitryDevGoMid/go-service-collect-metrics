package repository

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/models"
)

type RepositoryMetrics interface {
	ChangeMetrics() error
	GetGaugeMetricsAll() (*models.AllMetrics, error)
	GetCounterMetricsAll() (*models.AllMetrics, error)
	changePollCount() int64
}

type repositoryMetrics struct {
	metrics    *models.MemStorage
	metricsAll *models.AllMetrics
}

// Сonstructor
func NewRepositoryMetrics(metrics *models.MemStorage, metricsAll *models.AllMetrics) RepositoryMetrics {
	return &repositoryMetrics{metrics: metrics, metricsAll: metricsAll}
}

// Change all metrics
func (rp *repositoryMetrics) ChangeMetrics() error {

	fmt.Println("Запустили обновление метрик!")

	defer fmt.Println("Завершили обновление метрик!")

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	// Runtime metrics
	rp.metrics.Gauge.Alloc = float64(rtm.Alloc)
	rp.metrics.Gauge.BuckHashSys = float64(rtm.BuckHashSys)
	rp.metrics.Gauge.Frees = float64(rtm.Frees)
	rp.metrics.Gauge.GCCPUFraction = float64(rtm.GCCPUFraction)
	rp.metrics.Gauge.GCSys = float64(rtm.GCSys)
	rp.metrics.Gauge.HeapAlloc = float64(rtm.HeapAlloc)
	rp.metrics.Gauge.HeapIdle = float64(rtm.HeapIdle)
	rp.metrics.Gauge.HeapInuse = float64(rtm.HeapInuse)
	rp.metrics.Gauge.HeapObjects = float64(rtm.HeapObjects)
	rp.metrics.Gauge.HeapReleased = float64(rtm.HeapReleased)
	rp.metrics.Gauge.HeapSys = float64(rtm.HeapSys)
	rp.metrics.Gauge.LastGC = float64(rtm.LastGC)
	rp.metrics.Gauge.Lookups = float64(rtm.Lookups)
	rp.metrics.Gauge.MCacheInuse = float64(rtm.MCacheInuse)
	rp.metrics.Gauge.MCacheSys = float64(rtm.MCacheSys)
	rp.metrics.Gauge.MSpanInuse = float64(rtm.MSpanInuse)
	rp.metrics.Gauge.MSpanSys = float64(rtm.MSpanSys)
	rp.metrics.Gauge.Mallocs = float64(rtm.Mallocs)
	rp.metrics.Gauge.NextGC = float64(rtm.NextGC)
	rp.metrics.Gauge.NumForcedGC = float64(rtm.NumForcedGC)
	rp.metrics.Gauge.NumGC = float64(rtm.NumGC)
	rp.metrics.Gauge.OtherSys = float64(rtm.OtherSys)
	rp.metrics.Gauge.PauseTotalNs = float64(rtm.PauseTotalNs)
	rp.metrics.Gauge.StackInuse = float64(rtm.StackInuse)
	rp.metrics.Gauge.StackSys = float64(rtm.StackSys)
	rp.metrics.Gauge.Sys = float64(rtm.Sys)
	rp.metrics.Gauge.TotalAlloc = float64(rtm.TotalAlloc)

	// Custom metrics
	rp.metrics.Counter.PollCount = rp.changePollCount()
	rp.metrics.Gauge.RandomValue = randomValue()

	return nil
}

// Get random value
func randomValue() float64 {
	return (rand.Float64() * 5) + 5
}

// Count metrics PollCount
func (rp *repositoryMetrics) changePollCount() int64 {
	return rp.metrics.Counter.PollCount + 1
}

// Get all GAUGE metrics, use reflection for get Name and Value by Map
func (rp repositoryMetrics) GetGaugeMetricsAll() (*models.AllMetrics, error) {

	structMetricsReflection := reflect.ValueOf(&rp.metrics.Gauge).Elem()

	for i := 0; i < structMetricsReflection.NumField(); i++ {
		name := structMetricsReflection.Type().Field(i).Name
		value := structMetricsReflection.Field(i).Interface()
		rp.metricsAll.Gauge[name] = float64(value.(float64))
	}

	return rp.metricsAll, nil
}

// Get all COUNTER metrics, use reflection for get Name and Value by Map
func (rp repositoryMetrics) GetCounterMetricsAll() (*models.AllMetrics, error) {

	structMetricsReflection := reflect.ValueOf(&rp.metrics.Counter).Elem()

	for i := 0; i < structMetricsReflection.NumField(); i++ {
		name := structMetricsReflection.Type().Field(i).Name
		value := structMetricsReflection.Field(i).Interface()
		rp.metricsAll.Counter[name] = int64(value.(int64))
	}

	return rp.metricsAll, nil
}
