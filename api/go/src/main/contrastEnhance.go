package main
import (
	"../image"
	"fmt"
	"sync"
)


func ContrastEnhancement(fileName string, isLocal bool){
	fmt.Println("Enhancing contrast with...", fileName)

	wg := sync.WaitGroup{}
	wg.Add(3)
	
	contrastEnhancementFunc := func(fn interface{}) {
		defer wg.Done()
		switch fn.(type) {
		case func(string, bool):
			fn.(func(string,bool))(fileName, isLocal)
		}
	}

	// RGB
	go contrastEnhancementFunc(image.RGBHistogramEquilization)

	// YUV
	go contrastEnhancementFunc(image.YUVHistogramEquilization)

	// HSL
	go contrastEnhancementFunc(image.HSLHistogramEquilization)

}