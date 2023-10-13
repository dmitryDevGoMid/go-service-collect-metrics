package routes

import (
	"github.com/gin-gonic/gin"
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
