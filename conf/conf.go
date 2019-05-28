package conf

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/RudyChow/code-runner/app/utils"
)

var (
	configPath = "./config.toml"
	Cfg        *Config
)

type Config struct {
	Languages map[string]*info
	Container *container
	Http      *http
	Docker	  *docker
}

//语言配置
type info struct {
	Extension string
	Cmd       []string
	Images    map[string]string
}

//容器配置
type container struct {
	MaxExcuteTime       int64
	MaxLogLength        int
	TemFilePath         string
	ContainerNamePrefix string
	Limit               *containerLimit
}

type containerLimit struct {
	Memory    int64
	PidsLimit int64
	DiskQuota int64
	CPUShares int64
	CPUPeriod int64
	CPUQuota  int64
}

//http配置
type http struct {
	Port int
	Mode string
}

//docekr配置
type docker struct{
	ApiVersion string
}

func init() {
	initPath()
	//读取配置文件
	if _, err := toml.DecodeFile(configPath, &Cfg); err != nil {
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

func initPath() {
	envPath := os.Getenv("CODE_RUNNER_CONFIG_PATH")
	SetConfigPath(envPath)
}

func SetConfigPath(path string) {
	if path == "" {
		return
	}
	configPath = path
}
