package models

import (
	"testing"
)

//整个流程的基本测试
func TestLanguagesFlow(t *testing.T) {
	testLanguages := map[string]*ContainerOption{
		"php": &ContainerOption{
			Image:          "php:7.3-alpine",
			Cmd:            []string{"php", "main.php"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main.php",
			TargetFilePath: "/tmp/main.php",
		},
		"golang": &ContainerOption{
			Image:          "golang:1.12-alpine",
			Cmd:            []string{"go", "run", "main.go"},
			SourceFilePath: "/Users/rudy/go/src/github.com/RudyChow/code-runner/test/example/main.go",
			TargetFilePath: "/tmp/main.go",
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

func flow(containerOption *ContainerOption) (string, error) {

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

	s, err := DockerRunner.LogContainer(id)
	if err != nil {
		return "", err
	}

	err = DockerRunner.RemoveContainer(id)
	if err != nil {
		return "", err
	}

	return s, nil
}
