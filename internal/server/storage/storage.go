package storage

import "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"

func NewMemStorage() *models.MemStorage {

	mapGauge := make(map[string]float64)
	mapCounter := make(map[string]int64)

	memStorage := models.MemStorage{Gauge: mapGauge, Counter: mapCounter}

	return &memStorage
}
