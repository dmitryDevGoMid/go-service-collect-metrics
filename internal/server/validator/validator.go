package validator

import "errors"

var (
	ErrEmptyNameMetrics   = errors.New("name metrics cant be empty")
	ErrMetricsKeyNotFound = errors.New(`{"message": "metrics key not found"}`)
	ErrPingDataBase       = errors.New(`{"message": "database oing is fail"}`)
	ErrPingBatchDataBase  = errors.New(`{"message": "not save database batch"}`)
	ErrNotFoundType       = errors.New(`{"message": "not found type metric name"}`)
	ErrBadRequest         = errors.New(`{"message": "is bad request"}`)
)
