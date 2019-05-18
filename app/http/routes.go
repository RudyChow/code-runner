package http

import (
	"github.com/RudyChow/code-runner/app/http/controllers"
	"github.com/gin-gonic/gin"
)

func registerRouters(r *gin.Engine) {
	registerApiRoutes(r)
}

func registerApiRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.POST("/code", controllers.GetResult)
		v1.GET("/versions", controllers.GetVersions)
	}
}
