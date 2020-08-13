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

	contrastFuncs := [3]interface{}{RGBHistogramEquilization, YUVHistogramEquilization, HSLHistogramEquilization}
	for fn := range contrastFuncs {
		go contrastEnhancementFunc(fn)
	}

	wg.Wait()
}