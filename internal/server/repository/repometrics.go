package repository

import (
	"context"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
)

type MetricsRepository interface {
	GetMetricGauge(ctx context.Context, nameMetric string) (float64, error)
	GetMetricCounter(ctx context.Context, nameMetric string) (int64, error)
	UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error
	UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error
	GetAllMetrics(ctx context.Context) (*models.MemStorage, error)
	PingDatabase(ctx context.Context) error
	SaveMetricsBatch(ctx context.Context, metrics []unserialize.Metrics) error
}

// Декорируем каждый запрос к базе чтобы можно было повторно выполнять в случаи сбоя
type Decorator struct {
	IMetric MetricsRepository
}

// Кол-во попыток с интервалом между ними 1,2,3 секунды
var countAttempt [4]int = [4]int{0, 1, 2, 3}

func (d Decorator) GetMetricCounter(ctx context.Context, nameMetric string) (int64, error) {
	var err error
	var val int64
	for i := 0; i < len(countAttempt); i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		val, err = d.IMetric.GetMetricCounter(ctx, nameMetric)
		if err == nil {
			return val, nil
		} else {
			time.Sleep(time.Duration(countAttempt[i]) * time.Second)
		}
	}

	return 0, err
}

func (d Decorator) GetMetricGauge(ctx context.Context, nameMetric string) (float64, error) {
	var err error
	var val float64
	for i := 0; i < len(countAttempt); i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		val, err = d.IMetric.GetMetricGauge(ctx, nameMetric)
		if err == nil {
			return val, nil
		} else {
			time.Sleep(time.Duration(countAttempt[i]) * time.Second)
		}
	}

	return 0, err
}

// UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error
func (d Decorator) UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error {
	var err error
	for i := 0; i < len(countAttempt); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := d.IMetric.UpdateMetricGauge(ctx, nameMetric, value)
		if err == nil {
			return nil
		} else {
			time.Sleep(time.Duration(countAttempt[i]) * time.Second)
		}
	}

	return err
}

// UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error
func (d Decorator) UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error {
	var err error
	for i := 0; i < len(countAttempt); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := d.IMetric.UpdateMetricCounter(ctx, nameMetric, value)
		if err == nil {
			return nil
		} else {
			time.Sleep(time.Duration(countAttempt[i]) * time.Second)
		}
	}

	return err
}

// GetAllMetrics(ctx context.Context) (*models.MemStorage, error)
func (d Decorator) GetAllMetrics(ctx context.Context) (*models.MemStorage, error) {
	var err error
	for i := 0; i < len(countAttempt); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		val, err := d.IMetric.GetAllMetrics(ctx)
		if err == nil {
			return val, nil
		} else {
			time.Sleep(time.Duration(countAttempt[i]) * time.Second)
		}
	}

	return nil, err
}
