package routes

import (
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
)

// Нарезали роуты
func NewGinMetricsRoutesChange(metricsHandlers handlers.MetricsHandlers) []*GinRoute {
	return []*GinRoute{
		{
			// Update url counter metrics
			Path:    "/",
			Method:  "get",
			Handler: metricsHandlers.GetAllMetricsHTML,
		},
		{
			Path:    "/update/:type/:metric/:value", // Регистр нового
			Method:  "post",
			Handler: metricsHandlers.Update,
		},
		{
			Path:    "/value/:type/:metric", // Логин
			Method:  "post",
			Handler: metricsHandlers.Value,
		},
		{
			Path:    "/value/:type/:metric", // Логин
			Method:  "get",
			Handler: metricsHandlers.Value,
		},
	}
}

/*
Пример запроса к серверу:

Скопировать кодTEXT
POST /update/counter/someMetric/527 HTTP/1.1
Host: localhost:8080
Content-Length: 0
Content-Type: text/plain
Пример ответа от сервера:

Скопировать кодTEXT
HTTP/1.1 200 OK
Date: Tue, 21 Feb 2023 02:51:35 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8
*/
