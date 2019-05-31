package http

import (
	"github.com/RudyChow/code-runner/app/http/controllers"
	"github.com/gin-gonic/gin"
)

func registerRouters(r *gin.Engine) {
	api := r.Group("/api/code")
	api.POST("/", controllers.GetResult)
	api.GET("/versions", controllers.GetVersions)
}
