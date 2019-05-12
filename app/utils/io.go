package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

//将字符串写入文件
//@param string writeString 写入的字符串
//@param string fileName    写入的文件
func WriteFile(writeString string, fileName string) error {
	if fileName == "" {
		return errors.New("give me a file name pls?")
	}

	var d1 = []byte(writeString)
	err := ioutil.WriteFile(fileName, d1, 0666)
	return err
}

//获取项目当前目录
func GetCurrPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

//生成一个文件名
func GenerateRandomFileName(withPath string, withExtension string) string {
	nanoTime := time.Now().UnixNano()
	fileName := strconv.FormatInt(nanoTime, 10)
	return path.Join(withPath, fileName+withExtension)
}
