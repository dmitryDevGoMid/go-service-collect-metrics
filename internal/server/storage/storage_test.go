package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemStorage(t *testing.T) {
	metricsName := "test"
	gauge := float64(500)
	counter := int64(500)

	storage := NewMemStorage()
	assert.NotNil(t, storage)
	assert.Empty(t, storage.Counter)
	assert.Empty(t, storage.Gauge)

	storage.Counter[metricsName] = counter
	storage.Gauge[metricsName] = gauge

	assert.NotEmpty(t, storage.Counter)
	assert.NotEmpty(t, storage.Gauge)

	assert.Equal(t, storage.Counter[metricsName], counter)
	assert.Equal(t, storage.Gauge[metricsName], gauge)
}
