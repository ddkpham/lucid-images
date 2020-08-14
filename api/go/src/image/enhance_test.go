package image

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHSLHistogramEqualizationConcurrent(t *testing.T) {
	fmt.Println("HSL test")
	img := "boat.jpeg"
	start := time.Now()
	HSLHistogramEqualizationSerial(img, true)
	end := time.Now()
	dur1 := end.Sub(start)


	start = time.Now()
	HSLHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)

	assert.True(t, dur2 < dur1, "Concurrent HSL histogram equalization did not go as expected")
}

func TestYUVHistogramEqualizationConcurrent(t *testing.T) {
	fmt.Println("HSL test")
	img := "boat.jpeg"
	start := time.Now()
	YUVHistogramEqualizationSerial(img, true)
	end := time.Now()
	dur1 := end.Sub(start)


	start = time.Now()
	YUVHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)

	assert.True(t, dur2 < dur1, "Concurrent YUV histogram equalization did not go as expected")
}

func TestRGBHistogramEqualizationConcurrent(t *testing.T) {
	fmt.Println("HSL test")
	img := "boat.jpeg"
	start := time.Now()
	RGBHistogramEqualizationSerial(img, true)
	end := time.Now()
	dur1 := end.Sub(start)


	start = time.Now()
	RGBHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)

	assert.True(t, dur2 < dur1, "Concurrent RGB histogram equalization did not go as expected")
}