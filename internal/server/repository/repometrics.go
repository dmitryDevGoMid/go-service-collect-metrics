package repository

import (
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
)

type MetricsRepository interface {
	GetMetricGauge(nameMetric string) (float64, error)
	GetMetricCounter(nameMetric string) (int64, error)
	UpdateMetricGauge(nameMetric string, value float64) error
	UpdateMetricCounter(nameMetric string, value int64) error
	//GetAllMetrics() *models.MemStorage
	GetAllMetrics() (*models.MemStorage, error)
	PingDatabase() error
	SaveMetricsBatch(metrics []unserialize.Metrics) error
}
