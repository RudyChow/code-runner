package models

import (
	"sort"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/RudyChow/code-runner/app/utils"
	"github.com/RudyChow/code-runner/conf"
)

type Language struct {
	Language       string `form:"language" json:"language" xml:"language"  binding:"required"`
	Version        string `form:"version" json:"version" xml:"version"  binding:"required"`
	Code           string `form:"code" json:"code" xml:"code"  binding:"required"`
	SourceFilePath string
}

// CheckVersion : 检查版本是否存在
func (this *Language) CheckVersion() bool {
	_, ok := conf.Cfg.Languages[this.Language].Images[this.Version]
	return ok
}

// GetImage : 获取镜像名称
func (this *Language) GetImage() string {
	image, _ := conf.Cfg.Languages[this.Language].Images[this.Version]
	return image
}

// GetCmd : 获取cmd
func (this *Language) GetCmd() []string {
	cmd := conf.Cfg.Languages[this.Language].Cmd
	return cmd
}

// GetExtension : 获取后缀
func (this *Language) GetExtension() string {
	ext := conf.Cfg.Languages[this.Language].Extension
	return ext
}

// GetContainerOption : 获取容器参数
func (this *Language) GetContainerOption() *common.ContainerOption {
	option := &common.ContainerOption{}

	option.Image = this.GetImage()
	option.Cmd = this.GetCmd()
	option.SourceFilePath = this.SourceFilePath
	option.TargetFilePath = this.generateTargetFilePath()

	return option
}

//输出文件
func (this *Language) OutputFile() error {
	fileName := utils.GenerateRandomFileName(conf.Cfg.Container.TemFilePath, this.GetExtension())

	err := utils.WriteFile(this.Code, fileName)
	if err == nil {
		this.SourceFilePath = fileName
	}

	return err
}

//获取mount进入的目录
func (this *Language) generateTargetFilePath() string {
	return "/tmp/main" + this.GetExtension()
}

// GetAllSupportedImages : 获取配置中的所有镜像
func GetAllSupportedImages() map[string][]string {

	result := make(map[string][]string)

	for language, info := range conf.Cfg.Languages {
		var images sort.StringSlice
		for _, image := range info.Images {
			images = append(images, image)
		}
		sort.Sort(sort.Reverse(images))
		result[language] = images
	}

	return result
}

// GetAllSupportedVersions : 获取配置中的所有版本
func GetAllSupportedVersions() map[string][]string {

	result := make(map[string][]string)

	for language, info := range conf.Cfg.Languages {
		var versions sort.StringSlice
		for version, _ := range info.Images {
			versions = append(versions, version)
		}
		sort.Sort(sort.Reverse(versions))
		result[language] = versions
	}

	return result
}
