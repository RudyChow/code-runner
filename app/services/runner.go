package services

import (
	"errors"

	"github.com/RudyChow/code-runner/app/channels"
	"github.com/RudyChow/code-runner/app/models"
)

//从容器中获取结果
func GetResultFromDocker(l *models.Language) (string, error) {
	//检查这个版本是否支持
	exist := l.CheckVersion()
	if !exist {
		return "", errors.New("unsupports version")
	}

	//输出文件
	err := l.OutputFile()
	if err != nil {
		return "", err
	}

	//获取容器参数
	option := l.GetContainerOption()

	//获取结果
	id, res, err := models.DockerRunner.Run(option)
	if err != nil {
		return "", err
	}

	//删除容器
	channels.RemoveContainerChan <- id
	// models.DockerRunner.RemoveContainer(id)
	//删除文件
	channels.RemoveFileChan <- l.SourceFilePath

	return res, nil
}
