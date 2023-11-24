package wpool

import (
	"context"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
	"github.com/go-resty/resty/v2"
)

type Result struct {
	Response *resty.Response
	Err      error
}

type Job struct {
	ExceFn sandlers.SandlerMetrics
	Args   []string
}

func (j *Job) execution(ctx context.Context) Result {
	response, err := j.ExceFn.SendMetrics(j.Args)

	if err != nil {
		return Result{
			Response: response,
			Err:      err,
		}
	}

	return Result{
		Response: response,
		Err:      nil,
	}
}
