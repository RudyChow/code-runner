package schedules

import (
	"fmt"
	"time"

	"github.com/RudyChow/code-runner/conf"
)

func RunTikers() {
	cleaner := time.NewTicker(time.Second * time.Duration(conf.Cfg.Container.MaxExcuteTime))

	for {
		select {
		case <-cleaner.C:
			// todo
			//清理没用的东西
			fmt.Println("start cleaning")
		}
	}
}
