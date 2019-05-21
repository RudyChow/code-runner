package models

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/RudyChow/code-runner/app/common"
)

//整个流程的基本测试
func TestLanguagesFlow(t *testing.T) {
	testLanguages := map[string]*common.ContainerOption{
		"php": &common.ContainerOption{
			Image:          "php:7.3-alpine",
			Cmd:            []string{"php", "main.php"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main.php",
			TargetFilePath: "/tmp/main.php",
		},
		"golang": &common.ContainerOption{
			Image:          "golang:1.12-alpine",
			Cmd:            []string{"go", "run", "main.go"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main.go",
			TargetFilePath: "/tmp/main.go",
		},
		"python2": &common.ContainerOption{
			Image:          "python:2.7-alpine",
			Cmd:            []string{"python", "main.py"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main2.py",
			TargetFilePath: "/tmp/main.py",
		},
		"python3": &common.ContainerOption{
			Image:          "python:3.6-alpine",
			Cmd:            []string{"python", "main.py"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main3.py",
			TargetFilePath: "/tmp/main.py",
		},
	}

	for _, containerOption := range testLanguages {
		res, err := flow(containerOption)
		if err != nil {
			t.Error(err)
		}
		// if language != res {
		// 	mes := fmt.Sprintf("unexcepted result:test-%s,given-%s", language, res)
		// 	t.Error(errors.New(mes))
		// }
		t.Log(res)
	}

}

func TestGetContainers(t *testing.T) {
	runner := DockerRunner
	_, err := runner.GetContainers()
	if err != nil {
		t.Error(err)
	}
}

func TestCleanExpiredContainers(t *testing.T) {
	runner := DockerRunner
	err := runner.CleanExpiredContainers(0)
	if err != nil {
		t.Error(err)
	}
}

func TestGetImages(t *testing.T) {
	list, err := DockerRunner.GetImages()
	if err != nil {
		t.Error(err)
	}

	for _, image := range list {
		t.Log(image.RepoTags)
	}

	// t.Log(list)
}

func TestPullImages(t *testing.T) {
	ioReader, err := DockerRunner.PullImage("php:5.6-alpine")
	if err != nil {
		t.Error(err)
	}
	io.Copy(os.Stdout, ioReader)
}

func TestGetAllUnsupportedImages(t *testing.T) {
	list := GetAllSupportedImages()

	for _, image := range list {
		t.Log(image)
	}

	// t.Log(list)
}

func TestStatContainer(t *testing.T) {
	_, err := DockerRunner.StatContainer("1e637b916f15")
	if err != nil {
		t.FailNow()
	}
}

func TestLogContainer(t *testing.T) {
	log, err := DockerRunner.LogContainer("8814c6fed507")
	if err != nil {
		t.FailNow()
	}
	fmt.Printf("%+v", log)
}

func flow(containerOption *common.ContainerOption) (string, error) {

	id, err := DockerRunner.CreateContainer(containerOption)
	if err != nil {
		return "", err
	}

	err = DockerRunner.StartContainer(id)
	if err != nil {
		return "", err
	}

	err = DockerRunner.WaitContainer(id)
	if err != nil {
		return "", err
	}

	logs, err := DockerRunner.LogContainer(id)
	if err != nil {
		return "", err
	}

	err = DockerRunner.RemoveContainer(id)
	if err != nil {
		return "", err
	}

	return logs.Out, nil
}
