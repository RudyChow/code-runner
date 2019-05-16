package services

import (
	"errors"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/RudyChow/code-runner/app/models"
)

//从容器中获取结果
func GetResultFromDocker(l *models.Language) (*common.ContainerResult, error) {
	//检查这个版本是否支持
	exist := l.CheckVersion()
	if !exist {
		return nil, errors.New("unsupports version")
	}

	//输出文件
	err := l.OutputFile()
	if err != nil {
		return nil, err
	}

	//获取容器参数
	option := l.GetContainerOption()

	//获取结果
	containerResult, err := models.DockerRunner.Run(option)
	if err != nil {
		return nil, err
	}

	//删除容器
	common.RemoveContainerChan <- containerResult.ID
	// models.DockerRunner.RemoveContainer(containerResult.ID)
	//删除文件
	common.RemoveFileChan <- l.SourceFilePath

	return containerResult, nil
}
