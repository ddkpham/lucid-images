package image

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	//"image/color"
	_ "image/jpeg"
	//"image/png"
	"os"
)

func pixelVal(x, y int, rectangle image.Rectangle) uint {
	return uint(y * rectangle.Max.X + x)
}

func YUVHistogramEquilization(fileName string){
	img, err := os.Open(getFullPath(fileName))

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()
	fmt.Println("bounds: ", bounds)
	// convert image from rgb to yuv
	imgSize := bounds.Max.X * bounds.Max.Y
	y_img := make([]uint8, imgSize)
	cb_img := make([]uint8, imgSize)
	cr_img := make([]uint8, imgSize)
	a_img := make([]uint8, imgSize)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			// were only really interested in y for histogram equilization
			Y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
			y_img[pixelVal(x,y, bounds)] = Y
			cb_img[pixelVal(x,y, bounds)] = cb
			cr_img[pixelVal(x,y, bounds)] = cr
			a_img[pixelVal(x,y, bounds)] = uint8(a)
		}
	}

	// create histogram for y.
	y_hist := [256]uint32{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			y_hist[y_img[pixelVal(x,y, bounds)]]++
		}
	}

	yLUT := getLookUpTable(y_hist, imgSize)
	fmt.Println("y look up table: ", yLUT)

	//generate new contrast enhanced image with y Look up tables
	w, h := bounds.Max.X , bounds.Max.Y
	rect := image.Rect(0,0,w, h)
	newImg := image.NewRGBA64(rect)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r,g,b := color.YCbCrToRGB(yLUT[y_img[pixelVal(x,y, bounds)]], cb_img[pixelVal(x,y, bounds)], cr_img[pixelVal(x,y, bounds)])
			newImg.Set(x,y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a_img[pixelVal(x,y, bounds)],
			})
		}
	}
	f, err := os.Create("enhanced-YUV-" + fileName )
	if err != nil {
		panic(err)
	}

	defer f.Close()
	//err = png.Encode(f, newImg)
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}


}


// RGB contrast enhancement
func RGBHistogramEquilization(fileName string){
	img, err := os.Open(getFullPath(fileName))

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()
	fmt.Println("bounds: ", bounds)

	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	// https://golang.org/pkg/image/
	// generate histogram for each RGB channel.
	count := 0
	rHistogram, gHistogram, bHistogram := [256]uint32{}, [256]uint32{}, [256]uint32{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			count++
			r, g, b, a := decodedImg.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			rHistogram[r]++
			gHistogram[g]++
			bHistogram[b]++
		}
	}
	fmt.Println("count: ", count)

	// construct the Look Up Table by calculating the CDF
	imgSize := bounds.Max.Y * bounds.Max.X
	rLUT := getLookUpTable(rHistogram, imgSize)
	fmt.Println("look up table: ", rLUT)
	gLUT := getLookUpTable(gHistogram, imgSize)
	bLUT := getLookUpTable(bHistogram, imgSize)

	//generate new contrast enhanced image with Look up tables
	w, h := bounds.Max.X , bounds.Max.Y
	rect := image.Rect(0,0,w, h)
	newImg := image.NewRGBA64(rect)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			newImg.Set(x,y, color.RGBA{
				R: rLUT[r],
				G: gLUT[g],
				B: bLUT[b],
				A: uint8(a),
			})
		}
	}

	f, err := os.Create("enhanced-" + fileName )
	if err != nil {
		panic(err)
	}

	defer f.Close()
	//err = png.Encode(f, newImg)
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
	//// Print the results.

}


func getLookUpTable(histogram [256]uint32, imageSize int) [256]uint8 {
	// construct look up table by caluclating CDF
	fmt.Println("histogram: ", histogram)
	fmt.Println("imageSize: ", imageSize)
	sum := uint64(0)
	for _, val := range histogram {
		sum += uint64(val)
	}
	fmt.Println("sum ", sum)
	cdf := uint32(0)
	min := uint32(0)
	i := 0


	// find first non-zero value in histogram
	for {
		min = histogram[i]
		if min != 0 {
			break
		}
		i++
	}


	d := float64(imageSize) - float64(min)
	lut := [256]uint8{}
	lut_sum := uint8(0)
	for i := 0 ; i < 256 ; i++ {
		cdf += uint32(histogram[i])
		mappedValue := (float64(cdf) - float64(min))*255/d + 0.5
		scaledValue := uint8(mappedValue)
		lut[i] = scaledValue
		// trim off any values over 255 and under 0
		if lut[i] < 0 {
			lut[i] = 0
		}

		if lut[i] > 255 {
			lut[i] = 255
		}
		lut_sum += lut[i]
	}
	fmt.Println("lut sum: ", lut_sum)
	return lut
}
