package easyrunner

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/wpool"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
)

type Run interface {
	ChangeMetricsByTimeGopsUtil(ctx context.Context)
	ChangeMetricsByTime(ctx context.Context)
	SendMetricsByTime(ctx context.Context)
}

type run struct {
	sandlers          sandlers.SandlerMetrics
	workerPoolByDenis *wpool.WorkerPool
	cfg               *config.Config
}

func NewRunner(sandlers sandlers.SandlerMetrics, cfg *config.Config) Run {
	workerPoolByDenis := wpool.NewWorkerPool(cfg.Workers.LimitWorkers)

	return &run{sandlers: sandlers, cfg: cfg, workerPoolByDenis: workerPoolByDenis}
}

// Change metrics
func (r *run) ChangeMetricsByTime(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.cfg.Metrics.PollInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ChangeMetricsByTime -> Stop")
			return
		case <-ticker.C:
			r.sandlers.ChangeMetrics()
		}
	}
}

// Change metrics
func (r *run) ChangeMetricsByTimeGopsUtil(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.cfg.Metrics.PollInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ChangeMetricsByTimeGopsUtil -> Stop")
			return
		case <-ticker.C:
			r.sandlers.ChangeMetricsGopsUtil()
		}
	}
}

// Send metrics
func (r *run) SendMetricsByTime(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(r.cfg.Metrics.ReportInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("SendMetricsByTime -> Stop")
			return
		case <-ticker.C:
			r.SendRun(ctx)

		}
	}
}

// Запускаем Poll Workers
func (r *run) SendRun(ctx context.Context) {
	//Дергаем метрики из репозитария
	metrics := r.sandlers.GetMetricsListAndBatch()

	for k, v := range metrics.MetricsList {

		dataForSend := v

		sendData := []string{dataForSend}

		send := func() {
			r.sandlers.SendMetrics(sendData)
		}

		r.workerPoolByDenis.AddTask(send, k)
	}
}
