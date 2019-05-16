package schedules

import (
	"time"

	"github.com/RudyChow/code-runner/app/models"
	"github.com/RudyChow/code-runner/app/utils"
	"github.com/RudyChow/code-runner/conf"
)

func RunTikers() {
	cleaner := time.NewTicker(time.Second * time.Duration(conf.Cfg.Container.MaxExcuteTime))

	for {
		select {
		case <-cleaner.C:
			//gc
			go models.DockerRunner.CleanExpiredContainers(conf.Cfg.Container.MaxExcuteTime * 20)
			go utils.CleanExpiredTempFiles(conf.Cfg.Container.TemFilePath, conf.Cfg.Container.MaxExcuteTime*20)
		}
	}
}
