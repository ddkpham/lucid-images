package main
import (
	"../image"
	"fmt"
)

func ContrastEnhancement(fileName string){
	fmt.Println("Enhancing contrast with...", fileName)

	// RGB
	image.RGBHistogramEquilization(fileName)

	// YUV
	image.YUVHistogramEquilization(fileName)

	// HSL
	image.HSLHistogramEquilization(fileName)

}