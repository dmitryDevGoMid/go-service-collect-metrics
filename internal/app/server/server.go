package server

import (
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"
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

	metricsModel := storage.NewMemStorage()

	//Репозетарий
	metricsRepository := repository.NewMetricsRepository(metricsModel)

	//Обработчики
	metricsHandlers := handlers.NewMetricsHandlers(metricsRepository)

	//Роутинг
	metricsRotes := routes.NewGinMetricsRoutes(metricsHandlers)

	router := gin.Default()

	//Middleware Set Content TYPE
	router.Use(routes.WriteContentType())

	//Инициализируем роуты
	routes.InstallRouteGin(router, metricsRotes)

	// Сервер
	router.Run(fmt.Sprintf("%s", cfg.Server.Address))

}
