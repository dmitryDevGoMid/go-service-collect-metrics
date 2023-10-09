package agent

import (
	"fmt"
	"os"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/storage"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
)

func MonitorMetrics() {

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	storageMetrics := storage.NewMemStorage()

	storageAllMetrics := storage.NewAllMetricsStorage()

	repositoryMetrics := repository.NewRepositoryMetrics(storageMetrics, storageAllMetrics)

	sandlerMetrics := sandlers.NewMetricsSendler(repositoryMetrics)

	go sandlerMetrics.ChangeMetricsByTime(cfg)
	go sandlerMetrics.SendMetricsByTime(cfg)

	signalChannel := make(chan os.Signal, 1)

	<-signalChannel
}

func Run() {
	MonitorMetrics()
}
