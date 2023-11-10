package repository

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/models"
)

type RepositoryMetrics interface {
	ChangeMetrics() error
	GetGaugeMetricsAll() (*models.AllMetrics, error)
	GetCounterMetricsAll() (*models.AllMetrics, error)
	calcPollCount() int64
}

type repositoryMetrics struct {
	metrics    *models.MemStorage
	metricsAll *models.AllMetrics
	//Set mutex for gorutine metrics
	mutex *sync.Mutex
}

// Сonstructor пробрасываем две модели
func NewRepositoryMetrics(metrics *models.MemStorage, metricsAll *models.AllMetrics, mutex *sync.Mutex) RepositoryMetrics {
	return &repositoryMetrics{metrics: metrics, metricsAll: metricsAll, mutex: mutex}
}

// Change all metrics
func (rp *repositoryMetrics) ChangeMetrics() error {

	fmt.Println("Запустили обновление метрик!")

	defer func() {
		rp.mutex.Unlock()
		fmt.Println("Завершили обновление метрик!")
	}()

	rp.mutex.Lock()

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
	rp.metrics.Counter.PollCount = rp.calcPollCount()
	rp.metrics.Gauge.RandomValue = randomValue()

	return nil
}

// Get random value
func randomValue() float64 {
	return (rand.Float64() * 5) + 5
}

// Count metrics PollCount
func (rp *repositoryMetrics) calcPollCount() int64 {
	return rp.metrics.Counter.PollCount + 1
}

// Get all GAUGE metrics, use reflection for get Name and Value by Map
func (rp repositoryMetrics) GetGaugeMetricsAll() (*models.AllMetrics, error) {

	defer func() {
		rp.mutex.Unlock()
	}()

	rp.mutex.Lock()

	rp.metricsAll.Gauge["Alloc"] = rp.metrics.Gauge.Alloc
	rp.metricsAll.Gauge["BuckHashSys"] = rp.metrics.Gauge.BuckHashSys
	rp.metricsAll.Gauge["Frees"] = rp.metrics.Gauge.Frees
	rp.metricsAll.Gauge["GCCPUFraction"] = rp.metrics.Gauge.GCCPUFraction
	rp.metricsAll.Gauge["GCSys"] = rp.metrics.Gauge.GCSys
	rp.metricsAll.Gauge["HeapAlloc"] = rp.metrics.Gauge.HeapAlloc
	rp.metricsAll.Gauge["HeapIdle"] = rp.metrics.Gauge.HeapIdle
	rp.metricsAll.Gauge["HeapInuse"] = rp.metrics.Gauge.HeapInuse
	rp.metricsAll.Gauge["HeapObjects"] = rp.metrics.Gauge.HeapObjects
	rp.metricsAll.Gauge["HeapReleased"] = rp.metrics.Gauge.HeapReleased
	rp.metricsAll.Gauge["HeapSys"] = rp.metrics.Gauge.HeapSys
	rp.metricsAll.Gauge["LastGC"] = rp.metrics.Gauge.LastGC
	rp.metricsAll.Gauge["Lookups"] = rp.metrics.Gauge.Lookups
	rp.metricsAll.Gauge["MCacheInuse"] = rp.metrics.Gauge.MCacheInuse
	rp.metricsAll.Gauge["MCacheSys"] = rp.metrics.Gauge.MCacheSys
	rp.metricsAll.Gauge["MSpanInuse"] = rp.metrics.Gauge.MSpanInuse
	rp.metricsAll.Gauge["MSpanSys"] = rp.metrics.Gauge.MSpanSys
	rp.metricsAll.Gauge["Mallocs"] = rp.metrics.Gauge.Mallocs
	rp.metricsAll.Gauge["NextGC"] = rp.metrics.Gauge.NextGC
	rp.metricsAll.Gauge["NumForcedGC"] = rp.metrics.Gauge.NumForcedGC
	rp.metricsAll.Gauge["NumGC"] = rp.metrics.Gauge.NumGC
	rp.metricsAll.Gauge["OtherSys"] = rp.metrics.Gauge.OtherSys
	rp.metricsAll.Gauge["PauseTotalNs"] = rp.metrics.Gauge.PauseTotalNs
	rp.metricsAll.Gauge["StackInuse"] = rp.metrics.Gauge.StackInuse
	rp.metricsAll.Gauge["StackSys"] = rp.metrics.Gauge.StackSys
	rp.metricsAll.Gauge["Sys"] = rp.metrics.Gauge.Sys
	rp.metricsAll.Gauge["TotalAlloc"] = rp.metrics.Gauge.TotalAlloc
	rp.metricsAll.Gauge["RandomValue"] = rp.metrics.Gauge.RandomValue

	//fmt.Println(rp.metricsAll)

	//var metrics models.AllMetrics

	//metrics := *rp.metricsAll

	return rp.metricsAll, nil
	//return &metrics, nil
}

// Get all COUNTER metrics, use reflection for get Name and Value by Map
func (rp repositoryMetrics) GetCounterMetricsAll() (*models.AllMetrics, error) {

	defer func() {
		rp.mutex.Unlock()
	}()

	rp.mutex.Lock()

	rp.metricsAll.Counter["PollCount"] = rp.metrics.Counter.PollCount

	return rp.metricsAll, nil
}
