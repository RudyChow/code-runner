package http

import (
	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	r := gin.Default()
	registerRouters(r)
	r.Run("0.0.0.0:8080") // listen and serve on
}
