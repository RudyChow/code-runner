package main

import (
	"testing"

	"github.com/RudyChow/code-runner/app/models"
)

func TestMain(t *testing.T) {

}

func TestLanguage(t *testing.T) {
	l := &models.Language{
		Language: "php",
		Version:  "7.3",
		Code:     "<?php\necho 'ok';",
	}

	res := l.CheckVersion()
	if !res {
		t.Error("version not exist")
	}

	cmd := l.GetCmd()
	t.Log(cmd)

	ext := l.GetExtension()
	t.Log(ext)

	err := l.OutputFile()
	if err != nil {
		t.Error(err)
	}

	option := l.GetContainerOption()
	t.Log(option)

}

func TestTask(t *testing.T) {
	l := &models.Language{
		Language: "golang",
		Version:  "1.12",
		Code:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfor{fmt.Print(\"golang\")}\n}",
	}

	err := l.OutputFile()
	if err != nil {
		t.Error(err)
	}

	option := l.GetContainerOption()
	t.Log(option)

	containerResult, err := models.DockerRunner.Run(option)
	if err != nil {
		t.Error(err)
	}
	t.Log(containerResult.ID)
	t.Log(containerResult.Result)
}
