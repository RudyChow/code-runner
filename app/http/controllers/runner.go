package controllers

import (
	"net/http"

	"github.com/RudyChow/code-runner/app/models"
	"github.com/RudyChow/code-runner/app/services"
	"github.com/gin-gonic/gin"
)

//获取结果
func GetResult(c *gin.Context) {
	var l *models.Language

	if err := c.ShouldBind(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := services.GetResultFromDocker(l)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": s,
	})
}
