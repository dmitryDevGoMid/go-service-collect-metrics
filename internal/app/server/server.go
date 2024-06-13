package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config/db"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers/hgrpc"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/migration"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pb"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/cryptohashsha"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/logger"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/file"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mrepository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"
	"google.golang.org/grpc"

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

	//router := gin.Default()
	router := gin.New()

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

	hash256 := cryptohashsha.NewSha256(cfg)

	router.Use(routes.CheckHashSHA256Data(cfg, hash256))

	asme := asimencrypt.NewAsimEncrypt(cfg)

	errSetPrivateKey := asme.SetPrivateKey()

	if errSetPrivateKey != nil {
		log.Println("error set private key:", errSetPrivateKey)
	}

	router.Use(routes.AssimEncryptBody(cfg, asme))

	router.Use(routes.ToolsGroupPermission())

	//router.Use(routes.ToolsGroupPermission())

	//Инициализируем роуты
	routes.InstallRouteGin(router, metricsRotes)

	dbMigration := migration.NewMigration(dbConnection.DB())

	dbMigration.RunDrop(ctx)
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

	//Принимает данные по протоколу GRPC
	if cfg.TypeProtocolForSend.GetByGRPC {
		//Получаем обработчики gRPC
		handlerGRPC := hgrpc.NewGRPCHandlers(activeRepository)

		//Запускаем сервер gRPC
		go RunGRPCServer(handlerGRPC, cfg)

		/*if errGRPC != nil {
			log.Fatalf("not run serverGRPC: %s\n", errGRPC)
		} else {
			fmt.Println("run server GRPC")
		}*/
	}

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

func RunGRPCServer(server *hgrpc.ServerGRPC, cfg *config.Config) error {

	adress := cfg.ServerGRPC.AddressGRPC

	lis, err := net.Listen("tcp", adress)
	if err != nil {
		//return log.Fatalf("Failed to listen: %v", err)
		cfg.TypeProtocolForSend.GetByGRPC = false
		return err
	}

	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, server)

	// Set up graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("GRPC1")
	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")
		s.GracefulStop()
	}()

	fmt.Println("GRPC2")

	var errServer error

	//go func() {
	if errServer = s.Serve(lis); err != nil {
		//log.Fatalf("Failed to serve: %v", err)
		cfg.TypeProtocolForSend.GetByGRPC = false
		//return err
	}
	//}()

	if errServer != nil {
		return err
	}

	fmt.Println("GRPC3")

	return nil
}
