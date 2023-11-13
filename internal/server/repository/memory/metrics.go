package memory

import (
	"context"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/validator"
)

type metricsRepository struct {
	metrics *models.MemStorage
}

// Contruct
func NewMetricsRepository(metrics *models.MemStorage) repository.MetricsRepository {
	return &metricsRepository{metrics: metrics}
}

// Get matrics Gauge
func (mr *metricsRepository) GetMetricGauge(ctx context.Context, nameMetric string) (float64, error) {

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
func (mr *metricsRepository) GetMetricCounter(ctx context.Context, nameMetric string) (int64, error) {

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
func (mr *metricsRepository) UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error {

	if nameMetric == "" {
		return validator.ErrEmptyNameMetrics
	}
	mr.metrics.Gauge[nameMetric] = value

	return nil
}

// Upodate metrics Counter
func (mr *metricsRepository) UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error {

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
func (mr *metricsRepository) GetAllMetrics(ctx context.Context) (*models.MemStorage, error) {

	return mr.metrics, nil
}

func (mr *metricsRepository) PingDatabase(ctx context.Context) error {
	return validator.ErrPingDataBase
}

func (mr *metricsRepository) SaveMetricsBatch(ctx context.Context, metrics []unserialize.Metrics) error {
	for _, val := range metrics {
		if val.MType == "gauge" {
			err := mr.UpdateMetricGauge(ctx, val.ID, *val.Value)
			if err != nil {
				return err
			}
		}
		if val.MType == "counter" {
			err := mr.UpdateMetricCounter(ctx, val.ID, *val.Delta)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
