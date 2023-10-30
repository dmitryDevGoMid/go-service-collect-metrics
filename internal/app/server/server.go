package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/pkg/logger"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/routes"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//Инициализируем логгер
	appLogger := logger.NewAPILogger(cfg)

	appLogger.InitLogger()
	appLogger.Info("Start Service API Metrics")

	appLogger.Infof(
		"AppVersion: %s",
		cfg.Logger.Level,
	)

	metricsModel := storage.NewMemStorage()

	//Репозетарий
	metricsRepository := repository.NewMetricsRepository(metricsModel)

	//Обработчики
	metricsHandlers := handlers.NewMetricsHandlers(metricsRepository, cfg)

	//Роутинг
	metricsRotes := routes.NewGinMetricsRoutesChange(metricsHandlers)

	//router := gin.Default()
	router := gin.New()

	//Middleware Logger
	router.Use(routes.LoggerMiddleware(appLogger))

	//Middleware Set Content TYPE
	router.Use(routes.WriteContentType())

	//Middleware CORS
	//router.Use(routes.CORSMiddleware())

	//Инициализируем роуты
	routes.InstallRouteGin(router, metricsRotes)

	// Сервер
	//router.Run(cfg.Server.Address)

	// Line 27
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	log.Println("Shutdown Server ...")

	// Line 49
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Line 51
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
