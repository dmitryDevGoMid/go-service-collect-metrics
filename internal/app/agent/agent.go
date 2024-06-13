package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/middleware"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/cryptohashsha"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers/easyrunner"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/sandlers/hgrpc"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/storage"
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

	sha256 := cryptohashsha.NewSha256(cfg)

	//Шифрование тела сообщения
	asme := asimencrypt.NewAsimEncrypt(cfg)
	errSetPrivateKey := asme.SetPublicKey()
	if errSetPrivateKey != nil {
		log.Println("error set private key:", errSetPrivateKey)
	}

	clientMiddleware := middleware.NewClientMiddleware(client, cfg, sha256, asme)

	clientMiddleware.AssimEncryptBody()
	clientMiddleware.OnBeforeRequest()
	clientMiddleware.OnAfterResponse()

	storageMetrics := storage.NewMemStorage()

	storageAllMetrics := storage.NewAllMetricsStorage()

	repositoryMetrics := repository.NewRepositoryMetrics(storageMetrics, storageAllMetrics, &mutex)

	sandlerMetrics := sandlers.NewMetricsSendler(repositoryMetrics, client, ctx, cfg)

	sandlerByGRPC, err := hgrpc.GetClient(ctx, cancel, cfg)
	if err != nil {
		fmt.Println("Error getting client porotobuff metrics:", err)
	}

	fmt.Println(sandlerByGRPC)
	runSend := easyrunner.NewRunner(sandlerMetrics, sandlerByGRPC, cfg)

	//Решение: запускаем две горутины по сборке и одну с Poll Workers для отправки
	go runSend.ChangeMetricsByTime(ctx)
	go runSend.SendMetricsByTime(ctx)
	go runSend.ChangeMetricsByTimeGopsUtil(ctx)

	//Решение: запускаем две горутины по сборке и по отправке
	//go sandlerMetrics.ChangeMetricsByTime()
	//go sandlerMetrics.SendMetricsByTime()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel

	cancel()

	//Даем время для завершения всех горутин которые были запущены
	time.Sleep(time.Second * 1)

	log.Println("Shutdown Agent ...")
}

func Run() {
	MonitorMetricsRun()
}
