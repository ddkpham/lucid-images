package image

import (
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
	if err != nil {
		panic(err)
	}
	imgPath := workDir + imageDir + fileName
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