package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/middleware"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/storage"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
	"github.com/go-resty/resty/v2"
)

var mutex sync.Mutex

func MonitorMetricsRun() {

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	client := resty.New()

	clientMiddleware := middleware.NewClientMiddleware(client, cfg)

	clientMiddleware.OnBeforeRequest()
	clientMiddleware.OnAfterResponse()

	storageMetrics := storage.NewMemStorage()

	storageAllMetrics := storage.NewAllMetricsStorage()

	repositoryMetrics := repository.NewRepositoryMetrics(storageMetrics, storageAllMetrics, &mutex)

	sandlerMetrics := sandlers.NewMetricsSendler(repositoryMetrics, client, ctx, cfg)

	go sandlerMetrics.ChangeMetricsByTime()
	go sandlerMetrics.SendMetricsByTime()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel

	cancel()

	//Даем время для завершения всех горутин которые были запущены
	//time.Sleep(time.Second * 10)

	log.Println("Shutdown Agent ...")
}

func Run() {
	MonitorMetricsRun()
}
