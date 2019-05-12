package utils

import (
	"testing"
)

func TestWriteFile(t *testing.T) {
	code := "<?php\necho 'hello world';"
	err := WriteFile(code, "./main.php")
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateRandomFileName(t *testing.T) {
	path, _ := GetCurrPath()

	fileName := GenerateRandomFileName(path, ".go")

	t.Log(fileName)
}

func TestCleanExpiredTempFiles(t *testing.T) {
	CleanExpiredTempFiles("/Users/rudy/go/src/github.com/RudyChow/code-runner/tmp", 2)
}
