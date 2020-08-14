package image

import (
	"fmt"
	"sync"
)


func ContrastEnhancement(fileName string, isLocal bool){
	fmt.Println("Enhancing contrast with...", fileName)
	format, _ := GuessImageFormat(fileName, isLocal)
	fmt.Println("format: ", format)

	if !(format == "jpeg" || format == "png") {
		fmt.Println("file must be png or jpeg... ", format)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	contrastEnhancementFunc := func(fn interface{}) {
		defer wg.Done()
		switch fn.(type) {
		case func(string, bool):
			fn.(func(string,bool))(fileName, isLocal)
		}
	}

	// Serial implementation
	//contrastFuncs := [3]interface{}{RGBHistogramEqualizationSerial, YUVHistogramEqualizationSerial, HSLHistogramEqualizationSerial}
	//for _, fn := range contrastFuncs {
	//	go contrastEnhancementFunc(fn)
	//}

	contrastFuncs := [3]interface{}{
		RGBHistogramEqualizationConcurrent,
		YUVHistogramEqualizationConcurrent,
		HSLHistogramEqualizationConcurrent,
	}
	for _, fn := range contrastFuncs {
		go contrastEnhancementFunc(fn)
	}

	wg.Wait()
}