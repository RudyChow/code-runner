package http

import (
	"strconv"

	"github.com/RudyChow/code-runner/conf"
	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	r := gin.Default()
	registerRouters(r)
	gin.SetMode(conf.Cfg.Http.Mode)
	r.Run("0.0.0.0:" + strconv.Itoa(conf.Cfg.Http.Port)) // listen and serve on
}
