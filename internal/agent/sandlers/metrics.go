package sandlers

import (
	"fmt"
	"log"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"

	"github.com/go-resty/resty/v2"
)

type SandlerMetrics interface {
	ChangeMetricsByTime(cfg *config.Config)
	SendMetricsByTime(cfg *config.Config)
	SendMetrics(cfg *config.Config)
}

type sandlerMetrics struct {
	repository repository.RepositoryMetrics
}

func NewMetricsSendler(repository repository.RepositoryMetrics) SandlerMetrics {
	return &sandlerMetrics{repository: repository}
}

// Change metrics
func (rm *sandlerMetrics) ChangeMetricsByTime(cfg *config.Config) {
	secondChange := time.Duration(cfg.Metrics.PollInterval)
	for {
		// Run change metrics before sleep 2 seconds
		time.Sleep(secondChange * time.Second)
		rm.repository.ChangeMetrics()
	}
}

// Send metrics
func (rm *sandlerMetrics) SendMetricsByTime(cfg *config.Config) {
	secondSend := time.Duration(cfg.Metrics.ReportInterval)

	for {
		// Run change metrics before sleep 2 seconds
		time.Sleep(secondSend * time.Second)
		rm.SendMetrics(cfg)
	}
}

// Send metrics to server http://localhost:8080
func (rm *sandlerMetrics) SendMetrics(cfg *config.Config) {
	fmt.Printf("Запустили отправку метрик: http://%s\n", cfg.Server.Address)

	defer fmt.Printf("Завершили отправку метрик: http://%s\n", cfg.Server.Address)

	client := resty.New()

	metrics, _ := rm.repository.GetGaugeMetricsAll()

	//Send metrics GAUGE
	for key, val := range metrics.Gauge {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%v", cfg.Server.Address, key, val)
		_, err := client.R().Post(url)
		if err != nil {
			panic(err)
		}
		//fmt.Println(response)
	}

	metrics, _ = rm.repository.GetCounterMetricsAll()

	//Send metrics COUNTER
	for key, val := range metrics.Counter {
		url := fmt.Sprintf("http://%s/update/counter/%s/%v", cfg.Server.Address, key, val)
		_, err := client.R().Post(url)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(response)
	}
}
