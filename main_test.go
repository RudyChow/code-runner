package main

import (
	"testing"

	"github.com/RudyChow/code-runner/app/models"
)

func TestMain(t *testing.T) {

}

func TestLanguage(t *testing.T) {
	l := &models.Language{
		Name:    "php",
		Version: "7.3",
		Code:    "<?php\necho 'ok';",
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
		Name:    "php",
		Version: "7.3",
		Code:    "<?php\necho 'ok';",
	}

	err := l.OutputFile()
	if err != nil {
		t.Error(err)
	}

	option := l.GetContainerOption()
	t.Log(option)

	id, res, err := models.DockerRunner.Run(option)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
	t.Log(res)
}
