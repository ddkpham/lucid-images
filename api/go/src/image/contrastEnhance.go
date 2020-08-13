package image

import (
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
	go contrastEnhancementFunc(RGBHistogramEquilization)

	// YUV
	go contrastEnhancementFunc(YUVHistogramEquilization)

	// HSL
	go contrastEnhancementFunc(HSLHistogramEquilization)


	wg.Wait()

}