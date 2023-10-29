package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/storage"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
)

var mutex sync.Mutex

func MonitorMetricsRun() {

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	storageMetrics := storage.NewMemStorage()

	storageAllMetrics := storage.NewAllMetricsStorage()

	repositoryMetrics := repository.NewRepositoryMetrics(storageMetrics, storageAllMetrics, &mutex)

	sandlerMetrics := sandlers.NewMetricsSendler(repositoryMetrics)

	go sandlerMetrics.ChangeMetricsByTime(cfg)
	go sandlerMetrics.SendMetricsByTime(cfg)

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	log.Println("Shutdown Agent ...")
}

func Run() {
	MonitorMetricsRun()
}
