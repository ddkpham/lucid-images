package main
import (
	"../image"
	"fmt"
)

func ContrastEnhancement(fileName string, isLocal bool){
	fmt.Println("Enhancing contrast with...", fileName)

	// RGB
	image.RGBHistogramEquilization(fileName, isLocal)

	// YUV
	image.YUVHistogramEquilization(fileName, isLocal)

	// HSL
	image.HSLHistogramEquilization(fileName, isLocal)

}