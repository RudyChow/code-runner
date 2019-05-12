package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
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
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path, nil
}

//生成一个文件名
func GenerateRandomFileName(withPath string, withExtension string) string {
	nanoTime := time.Now().UnixNano()
	fileName := strconv.FormatInt(nanoTime, 10)
	return path.Join(withPath, fileName+withExtension)
}

//清理过期文件
func CleanExpiredTempFiles(temPath string, gap int64) {
	list, err := ioutil.ReadDir(temPath)
	if err != nil {
		log.Println(err)
		return
	}

	if len(list) == 0 {
		return
	}

	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		go func(name string, created int64) {
			defer wg.Done()
			//没过期就不删除
			if created+gap > time.Now().Unix() {
				return
			}

			fileName := path.Join(temPath, name)
			err = os.Remove(fileName)
			if err != nil {
				log.Println("failed deleting file", fileName, err)
			}
			log.Println("success deleting file", fileName)
		}(v.Name(), v.ModTime().Unix())
	}

	wg.Wait()
}
