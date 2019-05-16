package models

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/RudyChow/code-runner/app/utils"
	"github.com/RudyChow/code-runner/conf"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var DockerRunner *Runner

//https://docs.docker.com/engine/api/v1.39/#operation/ContainerCreate
type Runner struct {
	dockerClient *client.Client
	ctx          context.Context
}

func init() {
	runner, err := newRunner()
	if err != nil {
		panic(err)
	}
	DockerRunner = runner
}

//创建一个docker执行者
func newRunner() (*Runner, error) {
	runner := &Runner{}

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	runner.dockerClient = cli
	runner.ctx = context.Background()
	return runner, nil
}

//创建容器
func (this *Runner) CreateContainer(containerOption *common.ContainerOption) (string, error) {
	resp, err := this.dockerClient.ContainerCreate(this.ctx, &container.Config{
		Image:      containerOption.Image,
		Cmd:        containerOption.Cmd,
		WorkingDir: "/tmp",
	}, &container.HostConfig{
		Mounts: []mount.Mount{ //docker 容器目录挂在到宿主机目录
			mount.Mount{
				Type:   mount.TypeBind,
				Source: containerOption.SourceFilePath,
				Target: containerOption.TargetFilePath,
			},
		},
		Resources: container.Resources{
			Memory:    conf.Cfg.Container.Limit.Memory * 1024 * 1024,
			PidsLimit: conf.Cfg.Container.Limit.PidsLimit,
			DiskQuota: conf.Cfg.Container.Limit.DiskQuota * 1024 * 1024,
			CPUShares: conf.Cfg.Container.Limit.CPUShares,
			CPUPeriod: conf.Cfg.Container.Limit.CPUPeriod * 1000,
			CPUQuota:  conf.Cfg.Container.Limit.CPUQuota * 1000,
		},
	}, nil, conf.Cfg.Container.ContainerNamePrefix+"-"+utils.GenerateRandomFileName("", ""))
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

//开始容器
func (this *Runner) StartContainer(containerId string) error {
	if err := this.dockerClient.ContainerStart(this.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

//等待容器
func (this *Runner) WaitContainer(containerId string) error {
	statusCh, errCh := this.dockerClient.ContainerWait(this.ctx, containerId, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	case <-time.After(time.Second * time.Duration(conf.Cfg.Container.MaxExcuteTime)):
		return errors.New("task time out")
	}
	return nil
}

//获取容器记录
func (this *Runner) LogContainer(containerId string) (result string, err error) {
	out, err := this.dockerClient.ContainerLogs(this.ctx, containerId, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(out)
	if err != nil {
		return "", err
	}

	result = string(data[0:100])

	return
}

//删除容器
func (this *Runner) RemoveContainer(containerId string) error {
	err := this.dockerClient.ContainerRemove(this.ctx, containerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	return err
}

//获取容器列表
func (this *Runner) GetContainers() ([]types.Container, error) {
	list, err := this.dockerClient.ContainerList(this.ctx, types.ContainerListOptions{
		All:   true,
		Limit: -1,
	})
	return list, err
}

//获取容器资源统计
func (this *Runner) StatContainer(id string) (io.ReadCloser, error) {
	result, err := this.dockerClient.ContainerStats(this.ctx, id, true)

	return result.Body, err
}

//获取镜像列表
func (this *Runner) GetImages() ([]types.ImageSummary, error) {
	list, err := this.dockerClient.ImageList(this.ctx, types.ImageListOptions{
		All: true,
	})
	return list, err
}

//拉取镜像
func (this *Runner) PullImage(image string) (io.ReadCloser, error) {
	ioReader, err := this.dockerClient.ImagePull(this.ctx, image, types.ImagePullOptions{})
	return ioReader, err
}

//删除过期容器
func (this *Runner) CleanExpiredContainers(gap int64) error {
	list, err := this.GetContainers()
	if err != nil {
		return err
	}

	if len(list) > 0 {
		wg := &sync.WaitGroup{}
		for _, container := range list {
			wg.Add(1)

			go func(names []string, id string, created int64) {
				defer wg.Done()
				//如果没过期，则不删除
				if created+gap > time.Now().Unix() {
					return
				}

				isRunner := false
				for _, name := range names {
					if strings.Contains(name, conf.Cfg.Container.ContainerNamePrefix) {
						isRunner = true
						break
					}
				}
				//如果是，则删除
				if isRunner {
					err = this.RemoveContainer(id)
					if err != nil {
						log.Println("cannot delete container", id, "reason:", err)
					} else {
						log.Println("success deleting container", id)
					}
				}
			}(container.Names, container.ID, container.Created)

		}
		wg.Wait()
	}

	return err
}

//执行
func (this *Runner) Run(containerOption *common.ContainerOption) (*common.ContainerResult, error) {
	returnData := &common.ContainerResult{}

	//新建容器
	id, err := DockerRunner.CreateContainer(containerOption)
	if err != nil {
		return returnData, err
	}
	returnData.ID = id
	//获取状态
	statReader, err := DockerRunner.StatContainer(id)
	if err != nil {
		return returnData, err
	}
	//收集状态结果
	var (
		statResult []*types.StatsJSON
		startTime  = time.Now()
	)
	go func(reader io.ReadCloser) {

		decoder := json.NewDecoder(reader)
		for {
			var streamData *types.StatsJSON
			if err = decoder.Decode(&streamData); err != nil {
				reader.Close()
				return
			}
			statResult = append(statResult, streamData)

		}

	}(statReader)
	//运行容器
	err = DockerRunner.StartContainer(id)
	if err != nil {
		return returnData, err
	}

	//等待容器状态停止运行
	DockerRunner.WaitContainer(id)
	endTime := time.Now()

	//获取容器日志
	returnData.Result, err = DockerRunner.LogContainer(id)
	if err != nil {
		return returnData, err
	}
	returnData.Stats, returnData.ExecutionTime = utils.ParseStat(statResult)
	if returnData.ExecutionTime == 0 {
		returnData.ExecutionTime = endTime.Sub(startTime).Nanoseconds() / 1e6
	}

	return returnData, err
}
