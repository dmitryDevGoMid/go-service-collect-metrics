package storage

import "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/models"

func NewMemStorage() *models.MemStorage {

	var metrics models.MemStorage

	return &metrics
}

func NewAllMetricsStorage() *models.AllMetrics {

	mapGauge := make(map[string]float64)
	mapCounter := make(map[string]int64)

	metrics := models.AllMetrics{Gauge: mapGauge, Counter: mapCounter}

	return &metrics
}
