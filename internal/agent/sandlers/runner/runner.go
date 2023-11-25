package runner

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
	sandlers sandlers.SandlerMetrics
	wpool    wpool.WorkerStack
	cfg      *config.Config
}

func NewRunner(sandlers sandlers.SandlerMetrics, cfg *config.Config) Run {
	wpool := wpool.New(1)
	return &run{sandlers: sandlers, cfg: cfg, wpool: wpool}
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

// Читаем канал ответов от сервера
func (r *run) ListingResponseServer(ctx context.Context) {
	result := r.wpool.ListResults()
	for {
		select {
		case response, ok := <-result:
			_ = response
			_ = ok
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}
}

// Send metrics
func (r *run) SendMetricsByTime(ctx context.Context) {
	//Запускаем воркеры
	go r.wpool.WorkerRun(ctx)

	r.wpool.RunJobs <- true

	//Генерируем задачи для воркеров
	go r.wpool.GenerateJob(ctx)

	//Слушаем результаты выполнения
	go r.ListingResponseServer(ctx)

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

func (r *run) SendRun(ctx context.Context) {
	metrics := r.sandlers.GetMetricsListAndBatch()

	var listJobs []wpool.Job

	for _, v := range metrics.MetricsList {

		dataForSend := v

		listJobs = append(listJobs, wpool.Job{ExceFn: r.sandlers, Args: []string{dataForSend}})
	}

	r.wpool.ListObjJobs <- listJobs
}
