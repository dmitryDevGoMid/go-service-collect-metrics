package repository

import (
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/validator"
)

type MetricsRepository interface {
	GetMetricGauge(nameMetric string) (float64, error)
	GetMetricCounter(nameMetric string) (int64, error)
	UpdateMetricGauge(nameMetric string, value float64) error
	UpdateMetricCounter(nameMetric string, value int64) error
	GetAllMetrics() *models.MemStorage
}

type metricsRepository struct {
	metrics *models.MemStorage
}

// Contruct
func NewMetricsRepository(metrics *models.MemStorage) MetricsRepository {
	return &metricsRepository{
		metrics: metrics,
	}
}

// Get matrics Gauge
func (mr *metricsRepository) GetMetricGauge(nameMetric string) (float64, error) {
	if nameMetric == "" {
		return 0, validator.ErrEmptyNameMetrics
	}

	val, ok := mr.metrics.Gauge[nameMetric]
	if ok {
		return val, nil
	}

	return 0, validator.ErrMetricsKeyNotFound

}

// Get metrics Counter
func (mr *metricsRepository) GetMetricCounter(nameMetric string) (int64, error) {
	if nameMetric == "" {
		return 0, validator.ErrEmptyNameMetrics
	}

	val, ok := mr.metrics.Counter[nameMetric]
	if ok {
		return val, nil
	}

	return 0, validator.ErrMetricsKeyNotFound

}

// Upodate metrics Gauge
func (mr *metricsRepository) UpdateMetricGauge(nameMetric string, value float64) error {
	if nameMetric == "" {
		return validator.ErrEmptyNameMetrics
	}
	mr.metrics.Gauge[nameMetric] = value

	return nil
}

// Upodate metrics Counter
func (mr *metricsRepository) UpdateMetricCounter(nameMetric string, value int64) error {
	if nameMetric == "" {
		return validator.ErrEmptyNameMetrics
	}

	//Sum metrics Counter
	val, ok := mr.metrics.Counter[nameMetric]
	if ok {
		mr.metrics.Counter[nameMetric] = value + val
	} else {
		mr.metrics.Counter[nameMetric] = value
	}

	return nil
}

// Get All Metrics
func (mr *metricsRepository) GetAllMetrics() *models.MemStorage {
	return mr.metrics
}
