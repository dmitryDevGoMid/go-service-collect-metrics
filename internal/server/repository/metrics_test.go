package repository

import (
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"

	"github.com/stretchr/testify/assert"
)

func TestGetMetricGauge(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	type CheckError func(want float64, metricValue float64, err error)

	tests := []struct {
		name              string
		metricNameStorage string
		metricName        string
		metricValue       float64
		want              float64
		checkError        CheckError
	}{
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500.123,
			want:              700.123,
			checkError: func(want float64, metricValue float64, err error) {
				assert.NotEqual(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500.123,
			want:              500.123,
			checkError: func(want float64, metricValue float64, err error) {
				assert.Equal(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsNotExist",
			metricValue:       500.123,
			want:              500.123,
			checkError: func(want float64, metricValue float64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsModel.Gauge[tt.metricNameStorage] = tt.metricValue
			metricValue, err := metricsRepository.GetMetricGauge(tt.metricName)
			tt.checkError(tt.want, metricValue, err)
		})
	}

}

func TestGetMetricCounter(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	type CheckError func(want int64, metricValue int64, err error)

	tests := []struct {
		name              string
		metricNameStorage string
		metricName        string
		metricValue       int64
		want              int64
		checkError        CheckError
	}{
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500,
			want:              700,
			checkError: func(want int64, metricValue int64, err error) {
				assert.NotEqual(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500,
			want:              500,
			checkError: func(want int64, metricValue int64, err error) {
				assert.Equal(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsNotExist",
			metricValue:       500,
			want:              500,
			checkError: func(want int64, metricValue int64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricsModel.Counter[tt.metricNameStorage] = tt.metricValue
			metricValue, err := metricsRepository.GetMetricCounter(tt.metricName)
			tt.checkError(tt.want, metricValue, err)
		})
	}

}

func TestUpdateMetricCounter(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	type CheckError func(want int64, metricValue int64, err error)

	tests := []struct {
		name              string
		metricNameStorage string
		metricName        string
		metricValue       int64
		want              int64
		checkError        CheckError
	}{
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500,
			want:              700,
			checkError: func(want int64, metricValue int64, err error) {
				assert.NotEqual(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500,
			want:              1000,
			checkError: func(want int64, metricValue int64, err error) {
				assert.Equal(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsNotExist",
			metricValue:       500,
			want:              500,
			checkError: func(want int64, metricValue int64, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "",
			metricValue:       500,
			want:              500,
			checkError: func(want int64, metricValue int64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := metricsRepository.UpdateMetricCounter(tt.metricNameStorage, tt.metricValue)
			assert.NoError(t, err)
			metricValue, err := metricsRepository.GetMetricCounter(tt.metricName)
			tt.checkError(tt.want, metricValue, err)
		})
	}

}

func TestUpdateMetricGauge(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	type CheckError func(want float64, metricValue float64, err error)

	tests := []struct {
		name              string
		metricNameStorage string
		metricName        string
		metricValue       float64
		want              float64
		checkError        CheckError
	}{
		{
			name:              "metricExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsExist",
			metricValue:       500.957,
			want:              700.957,
			checkError: func(want float64, metricValue float64, err error) {
				assert.NotEqual(t, want, metricValue)
				assert.NoError(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "metricsNotExist",
			metricValue:       500.957,
			want:              500.957,
			checkError: func(want float64, metricValue float64, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:              "metricNotExists",
			metricNameStorage: "metricsExist",
			metricName:        "",
			metricValue:       500.957,
			want:              500.957,
			checkError: func(want float64, metricValue float64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := metricsRepository.UpdateMetricGauge(tt.metricNameStorage, tt.metricValue)
			assert.NoError(t, err)
			metricValue, err := metricsRepository.GetMetricGauge(tt.metricName)
			tt.checkError(tt.want, metricValue, err)
		})
	}

}

func TestGetAllMetrics(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsRepository := NewMetricsRepository(metricsModel)

	allMetrics := metricsRepository.GetAllMetrics()
	assert.NotNil(t, allMetrics)
	assert.NotEmpty(t, allMetrics)

}
