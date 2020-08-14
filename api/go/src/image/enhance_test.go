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

	fmt.Println("HSL serial: ", dur1)


	start = time.Now()
	HSLHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)

	fmt.Println("HSL concurrent: ", dur2)

	assert.True(t, dur2 < dur1, "Concurrent HSL histogram equalization did not go as expected")
}

func TestYUVHistogramEqualizationConcurrent(t *testing.T) {
	fmt.Println("YUV test")
	img := "boat.jpeg"
	start := time.Now()
	YUVHistogramEqualizationSerial(img, true)
	end := time.Now()
	dur1 := end.Sub(start)
	fmt.Println("YUV serial: ", dur1)


	start = time.Now()
	YUVHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)
	fmt.Println("YUV concurrent: ", dur2)

	assert.True(t, dur2 < dur1, "Concurrent YUV histogram equalization did not go as expected")
}

func TestRGBHistogramEqualizationConcurrent(t *testing.T) {
	fmt.Println("RGB test")
	img := "boat.jpeg"
	start := time.Now()
	RGBHistogramEqualizationSerial(img, true)
	end := time.Now()
	dur1 := end.Sub(start)

	fmt.Println("RGB serial: ", dur1)


	start = time.Now()
	RGBHistogramEqualizationConcurrent(img, true)
	end = time.Now()
	dur2 := end.Sub(start)

	fmt.Println("RGB concurrent: ", dur2)

	assert.True(t, dur2 < dur1, "Concurrent RGB histogram equalization did not go as expected")
}