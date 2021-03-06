package schedules

import (
	"log"
	"os"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/RudyChow/code-runner/app/models"
)

//消费管道
func RunChan() {
	for {
		select {
		//删除容器
		case id := <-common.RemoveContainerChan:
			err := models.DockerRunner.RemoveContainer(id)
			if err != nil {
				log.Println("failed deleting container", id, err)
			} else {
				log.Println("success deleting container", id)
			}
		//删除文件
		case fileName := <-common.RemoveFileChan:
			err := os.Remove(fileName)
			if err != nil {
				log.Println("failed deleting file", fileName, err)
			} else {
				log.Println("success deleting file", fileName)
			}
		}
	}
}
