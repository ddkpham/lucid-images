package image

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	_ "golang.org/x/image/webp"
)

const VagrantImageDir string = "/home/vagrant/project/client/public/"
var VagrantGoPath string = "/home/vagrant/project/api/go/src/main"

func getFullPath(fileName string, isLocal bool) string {
	// imageDir := "/src/image/"
	workDir, err := os.Getwd()
	rootDir := strings.Split(workDir, "src")[0]
	imgDir := rootDir + "/src/image/"

	if err != nil {
		panic(err)
	}
	imgPath := imgDir + fileName // locally
	if !isLocal { // in VM
		imgPath = VagrantImageDir + fileName
	}
	fmt.Println("full image path: ", imgPath)
	// for when rabbit mq fails but still needs to process jobs
	//imgPath = "/home/vagrant/project/api/go/src/image/chain.png"
	return imgPath
}

// Guess image format from gif/jpeg/png/webp
func GuessImageFormat(fileName string, isLocal bool) (format string, err error) {
	fullPath := getFullPath(fileName, isLocal)
	img , err := os.Open(fullPath)
	_, format, err = image.DecodeConfig(img)
	return format, err
}
