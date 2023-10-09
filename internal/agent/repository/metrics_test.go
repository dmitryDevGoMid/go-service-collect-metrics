package repository

import (
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/storage"

	"github.com/stretchr/testify/assert"
)

func getStorages() (*models.AllMetrics, *models.MemStorage) {
	TestStorageMetrics := storage.NewMemStorage()

	TestStorageAllMetrics := storage.NewAllMetricsStorage()

	return TestStorageAllMetrics, TestStorageMetrics
}

func TestRepositoryMetrics(t *testing.T) {
	TestStorageAllMetrics, TestStorageMetrics := getStorages()
	assert.NotNil(t, TestStorageMetrics)
	assert.NotNil(t, TestStorageAllMetrics)

	TestRepositoryMetrics := NewRepositoryMetrics(TestStorageMetrics, TestStorageAllMetrics)
	assert.NotNil(t, TestRepositoryMetrics)

	TestRepositoryMetrics.ChangeMetrics()
	assert.NotNil(t, TestStorageAllMetrics)

	metricsCounterAll, err := TestRepositoryMetrics.GetCounterMetricsAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, metricsCounterAll)

	metricsGaugeAll, err := TestRepositoryMetrics.GetGaugeMetricsAll()
	assert.NoError(t, err)
	assert.NotEmpty(t, metricsGaugeAll)
}

func TestRandomValue(t *testing.T) {
	randomTest := randomValue()
	random := randomValue()

	assert.NotEmpty(t, random)
	assert.NotEqual(t, randomTest, random)
}
