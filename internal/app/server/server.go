package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config/db"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/migration"
	hashsha256 "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/cryptoSha256"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/logger"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/file"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mrepository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"

	repositoryDb "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/db"
	repositoryM "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/memory"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/routes"

	"github.com/gin-gonic/gin"
)

func Run() {

	ctx, cancel := context.WithCancel(context.Background())

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

	//var metricsRepository repository.MetricsRepository

	dbConnection := db.NewConnection(cfg)
	//err = dbConnection.Ping()

	metricsModel := storage.NewMemStorage()

	metricsRepositoryLocal := repositoryM.NewMetricsRepository(metricsModel)

	metricsRepositoryDB := repositoryDb.NewMetricsRepository(dbConnection.DB())

	managerRepository := mrepository.NewMamangerRepository(
		metricsRepositoryDB,
		metricsRepositoryLocal,
		dbConnection)

	activeRepository := managerRepository.GetRepositoryActive()

	//Обработчики
	metricsHandlers := handlers.NewMetricsHandlers(managerRepository, activeRepository, cfg)

	//Роутинг
	metricsRotes := routes.NewGinMetricsRoutesChange(metricsHandlers)

	router := gin.Default()
	//router := gin.New()

	// Работаем с временным файлом для сохранения данных из сервера
	workFile := file.NewWorkFile(metricsRepositoryLocal, cfg, ctx)

	// Запускаем процесс чтения и записи
	workFile.RunWorker(ctx)

	router.Use(routes.SaveFileToDisk(cfg, workFile))

	//Middleware Logger
	router.Use(routes.LoggerMiddleware(appLogger))

	//Middleware Set Content TYPE
	router.Use(routes.WriteContentType())

	//Middleware CORS
	router.Use(routes.CORSMiddleware())

	router.Use(routes.DecompressMiddleware())

	hash256 := hashsha256.NewSha256(cfg)

	router.Use(routes.CheckHashSHA256Data(cfg, hash256))

	//router.Use(routes.ToolsGroupPermission())

	//Инициализируем роуты
	routes.InstallRouteGin(router, metricsRotes)

	dbMigration := migration.NewMigration(dbConnection.DB())

	//dbMigration.RunDrop(ctx)
	dbMigration.RunCreate(ctx)

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

	cancel()

	dbConnection.Close()

	// Line 51
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
