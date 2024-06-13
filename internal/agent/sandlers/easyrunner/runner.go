package easyrunner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pb"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/serialize"
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
	sandlersByGRPC    pb.MetricsServiceClient
	workerPoolByDenis *wpool.WorkerPool
	cfg               *config.Config
}

func NewRunner(sandlers sandlers.SandlerMetrics, sandlersByGRPC pb.MetricsServiceClient, cfg *config.Config) Run {
	workerPoolByDenis := wpool.NewWorkerPool(cfg.Workers.LimitWorkers)

	return &run{sandlers: sandlers, sandlersByGRPC: sandlersByGRPC, cfg: cfg, workerPoolByDenis: workerPoolByDenis}
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

			if r.cfg.TypeProtocolForSend.SendByGRPC {
				r.SendByGRPC(ctx)
			} else {
				r.SendRun(ctx)
			}

		}
	}
}

// Запускаем Poll Workers
func (r *run) SendByGRPC(ctx context.Context) {
	//Дергаем метрики из репозитария
	metrics := r.sandlers.GetMetricsListAndBatch()

	var serializeMetrics serialize.Metrics

	for _, v := range metrics.MetricsList {

		dataForSend := v

		//sendData := []string{dataForSend}
		err := json.Unmarshal([]byte(dataForSend), &serializeMetrics)
		if err != nil {
			fmt.Println("unmarshal json err:", err)
		}

		var mr pb.MetricsRequest
		mr.Id = serializeMetrics.ID
		mr.Mtype = serializeMetrics.MType
		if mr.Mtype == "gauge" {
			mr.Value = *serializeMetrics.Value
		} else {
			mr.Delta = *serializeMetrics.Delta
		}

		fmt.Println("METRICS:", &mr)

		// Send MetricsRequest message to server

		send := func() {
			_, errMetrics := r.sandlersByGRPC.Metrics(ctx, &mr)
			if errMetrics != nil {
				//log.Fatalf("Failed to send metrics by GRPC: %v", err)
				fmt.Println("Failed to send metrics by GRPC:", errMetrics)
				return
			}
			//r.sandlers.SendMetrics(sendData)
		}
		send()

		//r.workerPoolByDenis.AddTask(send, k)
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
