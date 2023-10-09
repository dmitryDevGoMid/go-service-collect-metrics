package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type GinRoute struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}

// Set route by slice of struct GinRoute
func InstallRouteGin(router *gin.Engine, routeList []*GinRoute) {
	for _, route := range routeList {
		switch method := route.Method; method {
		case "post":
			router.POST(route.Path, route.Handler)
		default:
			router.GET(route.Path, route.Handler)
		}
	}
}

// Set politics CORS
func WithCORS(router *mux.Router) http.Handler {
	//Типы зоголовков
	//headers := handlers.AllowedHeaders([]string{"X-Requested-with", "Content-Type", "Accept", "Autorization"})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	//Разрешаем доступ с любого источника
	origins := handlers.AllowedOrigins([]string{"*"})
	//Доступные методы все основные
	//methods := handlers.AllowedMethods([]string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodGet})
	methods := handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet})

	return handlers.CORS(headers, methods, origins)(router)
}
