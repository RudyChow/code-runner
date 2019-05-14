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
	"fmt"
	"io"
	"os"

	"github.com/RudyChow/code-runner/app/models"
	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getImages()
		if err != nil {
			cmd.Println(err)
			os.Exit(-1)
		}

		formatPrintImages(result)

		var fix string
		cmd.Println("\nwant to download the nonexistent images?(y/n)")
		fmt.Scan(&fix)

		if fix == "n" {
			os.Exit(0)
		}

		for _, output := range result {
			//存在就不下载了
			if output.IsExist {
				continue
			}
			downloadImage(output.Image)
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
	fmt.Println("image\t", "exist")
	for _, output := range images {
		fmt.Println(output.Image+"\t", output.IsExist)
	}
}

//下载镜像
func downloadImage(image string) error {
	reader, err := models.DockerRunner.PullImage(image)
	if err != nil {
		return err
	}

	fmt.Println("start downloading image", image)

	io.Copy(os.Stdout, reader)

	fmt.Println("finish downloading image", image)

	return nil
}

type output struct {
	Image   string
	IsExist bool
}
