package models

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	"github.com/RudyChow/code-runner/conf"

	"github.com/RudyChow/code-runner/app/utils"

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
func (this *Runner) CreateContainer(containerOption *ContainerOption) (string, error) {
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
	}, nil, "code-runner-"+utils.GenerateRandomFileName("", ""))
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

	ticker := time.NewTicker(time.Second * time.Duration(conf.Cfg.Container.MaxExcuteTime))
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	case <-ticker.C:
		ticker.Stop()
		this.RemoveContainer(containerId)
		return errors.New("container time out")
	}
	return nil
}

//获取容器记录
func (this *Runner) LogContainer(containerId string) (string, error) {
	out, err := this.dockerClient.ContainerLogs(this.ctx, containerId, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(data), nil
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

//执行
func (this *Runner) Run(containerOption *ContainerOption) (string, string, error) {
	id, err := DockerRunner.CreateContainer(containerOption)
	if err != nil {
		return "", "", err
	}

	err = DockerRunner.StartContainer(id)
	if err != nil {
		return id, "", err
	}

	err = DockerRunner.WaitContainer(id)
	if err != nil {
		return id, "", err
	}

	s, err := DockerRunner.LogContainer(id)
	if err != nil {
		return id, "", err
	}

	return id, s, err
}