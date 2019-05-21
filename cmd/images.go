// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/RudyChow/code-runner/app/models"
	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "显示你当前支持的镜像状态以及自动下载缺失的镜像",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getImages()
		if err != nil {
			cmd.Println(err)
			os.Exit(-1)
		}

		formatPrintImages(result)

		var fix string
		fmt.Println("\ndo u want to download the nonexistent images?(y/N)")
		fmt.Scan(&fix)

		if fix == "y" {
			for _, output := range result {
				//存在就不下载了
				if output.IsExist {
					continue
				}
				downloadImage(output.Image)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
}

//获取当前配置的镜像是否已经拉取到本地
func getImages() ([]*output, error) {
	var result []*output

	supportedImages := models.GetAllSupportedImages()

	if len(supportedImages) == 0 {
		return result, nil
	}

	dockerImages, err := models.DockerRunner.GetImages()
	if err != nil {
		return result, err
	}
	for _, images := range supportedImages {
		for _, image := range images {
			//判断docker中是否有该镜像
			output := &output{
				Image:   image,
				IsExist: false,
			}
			for _, dockerImage := range dockerImages {
				for _, tag := range dockerImage.RepoTags {
					if image == tag {
						output.IsExist = true
						goto Next
					}
				}
			}
		Next:
			result = append(result, output)
		}
	}
	return result, nil
}

//格式化输出
func formatPrintImages(images []*output) {
	fmt.Printf("|%20s|%20s|\n", "image", "exist")
	for _, output := range images {
		fmt.Printf("|%20s|%20t|\n", output.Image, output.IsExist)
	}
}

//下载镜像
func downloadImage(image string) error {
	reader, err := models.DockerRunner.PullImage(image)
	if err != nil {
		return err
	}

	var result *pullStat
	// var result *map[string]interface{}
	decoder := json.NewDecoder(reader)

	fmt.Printf("start downloading %s\n", image)

	// io.Copy(os.Stdout, reader)
	for {
		if err := decoder.Decode(&result); err != nil {
			//读完就下一步
			if err == io.EOF {
				break
			}
			//有错就报错
			fmt.Printf("\nfailed downloading %s:%s\n", image, err)
			return err
		}
		// fmt.Printf("\r-%150s", "")
		fmt.Printf("\r%s:%s", result.Status, result.Progress)
	}

	fmt.Printf("\nsuccess downloading %s\n", image)

	return nil
}

type output struct {
	Image   string
	IsExist bool
}

type pullStat struct {
	Status   string
	Id       string
	Progress string
}
