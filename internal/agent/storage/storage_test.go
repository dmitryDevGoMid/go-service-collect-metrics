package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAllMetricsStorage(t *testing.T) {
	storage := NewAllMetricsStorage()
	assert.NotNil(t, storage)
}

func TestNewMemStorage(t *testing.T) {
	storage := NewMemStorage()
	assert.NotNil(t, storage)
}
