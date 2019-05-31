package controllers

import (
	"net/http"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/RudyChow/code-runner/app/models"
	"github.com/RudyChow/code-runner/app/services"
	"github.com/gin-gonic/gin"
)

//获取结果
func GetResult(c *gin.Context) {
	var l models.Language

	if err := c.ShouldBind(&l); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := services.GetResultFromDocker(&l)
	if err != nil {
		c.JSON(http.StatusForbidden, common.ApiResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.ApiResponse{
		Data: result,
	})
}

//获取支持的语言以及版本
func GetVersions(c *gin.Context) {
	result := models.GetAllSupportedVersions()
	c.JSON(http.StatusOK, common.ApiResponse{
		Data: result,
	})
}
