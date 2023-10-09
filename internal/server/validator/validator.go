package validator

import "errors"

var (
	ErrEmptyNameMetrics   = errors.New("name metrics cant be empty")
	ErrMetricsKeyNotFound = errors.New(`{"message": "metrics key not found"}`)
)
