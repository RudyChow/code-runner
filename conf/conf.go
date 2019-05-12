package conf

import (
	"log"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/RudyChow/code-runner/app/utils"
)

var Cfg *Config

type Config struct {
	Languages map[string]*info
	Container *container
}

type info struct {
	Extension string
	Cmd       []string
	Images    map[string]string
}

type container struct {
	MaxExcuteTime int64
	TemFilePath   string
}

func init() {
	//读取配置文件
	if _, err := toml.DecodeFile("./config.toml", &Cfg); err != nil {
		log.Panic(err)
	}
	//设置目录
	if Cfg.Container.TemFilePath == "" {
		currPath, err := utils.GetCurrPath()
		if err != nil {
			log.Panic(err)
		}
		Cfg.Container.TemFilePath = path.Join(currPath, "tmp")
	}

}
