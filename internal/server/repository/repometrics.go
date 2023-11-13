package repository

import (
	"context"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
)

type MetricsRepository interface {
	GetMetricGauge(ctx context.Context, nameMetric string) (float64, error)
	GetMetricCounter(ctx context.Context, nameMetric string) (int64, error)
	UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error
	UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error
	//GetAllMetrics() *models.MemStorage
	GetAllMetrics(ctx context.Context) (*models.MemStorage, error)
	PingDatabase(ctx context.Context) error
	SaveMetricsBatch(ctx context.Context, metrics []unserialize.Metrics) error
}
