package repository

import (
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricGauge(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	tests := []struct {
		name              string
		metricNameStorage string
		metricName        string
		metricValue       float64
		want              float64
	}{
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500.123,
			want:              500.123,
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsNotExist",
			metricValue:       500.123,
			want:              500.123,
		},
		{
			name:              "metricValueEqual",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500.123,
			want:              500.123,
		},
		{
			name:              "metricValueNotEqual",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       700.234,
			want:              500.123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "metricExists" {
				metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
				_, err := metricsRepository.GetMetricGauge(tt.metricName)
				assert.NoError(t, err)
			}
			if tt.name == "metricNotExists" {
				metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
				_, err := metricsRepository.GetMetricGauge(tt.metricName)
				assert.Error(t, err)
			}
			if tt.name == "metricNotExists" {
				metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
				_, err := metricsRepository.GetMetricGauge(tt.metricName)
				assert.Error(t, err)
			}
			if tt.name == "metricValueEqual" {
				metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
				val, err := metricsRepository.GetMetricGauge(tt.metricName)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, val)
			}
			if tt.name == "metricValueNotEqual" {
				metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
				val, err := metricsRepository.GetMetricGauge(tt.metricName)
				assert.NoError(t, err)
				assert.NotEqual(t, tt.want, val)
			}
		})
	}

}
