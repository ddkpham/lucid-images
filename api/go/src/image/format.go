package image

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"os"

	_ "golang.org/x/image/webp"
)

func getFullPath(fileName string) string {
	imageDir := "/src/image/"
	workDir, err := os.Getwd()
	fmt.Println("workDir: ", workDir)
	if err != nil {
		panic(err)
	}
	imgPath := workDir + imageDir + fileName // locally
	if workDir == "/home/vagrant/project/api/go/src/main" { // in VM
		imgPath = "/home/vagrant/project/api/go/src/image/" + fileName
	}
	// for when rabbit mq fails but still needs to process jobs
	//imgPath = "/home/vagrant/project/api/go/src/image/chain.png"
	return imgPath
}

// Guess image format from gif/jpeg/png/webp
func GuessImageFormat(fileName string) (format string, err error) {
	img , err := os.Open(getFullPath(fileName))
	_, format, err = image.DecodeConfig(img)
	return format, err
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(fileName string) string {
	format, _ := GuessImageFormat(fileName)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}