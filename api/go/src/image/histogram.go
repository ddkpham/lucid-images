package image

import (
	"image"
	"image/color"
	"image/jpeg"
	"sync"

	//"image/color"
	_ "image/jpeg"
	//"image/png"
	"os"
)


// Feel free to change this.
var numThreads int = 4

func pixelVal(x, y int, rectangle image.Rectangle) uint {
	return uint(y * rectangle.Max.X + x)
}


func maximum(r,g,b float32) float32 {
	max := r
	if g > max {
		max = g
	}

	if b > max {
		max = b
	}
	return max
}

func minimum(r,g,b float32) float32 {
	min := r
	if g < min {
		min = g
	}

	if b < min {
		min = b
	}
	return min
}

func rgb2hsl(r,g,b uint8) (h,s float32 ,l uint8) {
	hTemp, sTemp, lTemp := float32(0), float32(0), float32(0)
	// convert RGB to range [ 0, 1 ]
	rScaled := float32(r) / 255
	gScaled := float32(g) / 255
	bScaled := float32(b) / 255

	// find min & max of the rgb values
	min := minimum(rScaled, gScaled, bScaled)
	max := maximum(rScaled, gScaled, bScaled)

	deltaMax := max - min
	lTemp = (max + min) / 2
	if deltaMax == 0 { // Gray value. There is no chroma
		hTemp = 0
		sTemp = 0
	} else {
		if lTemp < 0.5 {
			sTemp = deltaMax / (min + max)
		} else {
			sTemp = deltaMax / (2 - (min + max))
		}

		deltaConversion := func(val float32) float32 {
			return (((max - val) / 6) + (deltaMax/2)) / deltaMax
		}

		rDelta := deltaConversion(rScaled)
		gDelta := deltaConversion(gScaled)
		bDelta := deltaConversion(bScaled)

		if rScaled == max {
			hTemp = bDelta - gDelta
		} else {
			if gScaled == max {
				hTemp = (1.0/3.0) + rDelta - bDelta
			} else {
				hTemp = (2.0/3.0) + gDelta - rDelta
			}
		}
	}

	if hTemp < 0 {
		hTemp+=1
	}

	if hTemp > 1 {
		hTemp-=1
	}

	h = hTemp
	s = sTemp
	l = uint8(lTemp * 255)
	return h,s,l
}

func hsl2rgb(h,s float32, l uint8) (r,g,b uint8) {
	var_1, var_2 := float32(0), float32(0)
	hTemp := h
	sTemp := s
	lTemp := float32(l) / 255

	if s == 0 {
		r = uint8(lTemp * 255)
		g = uint8(lTemp * 255)
		b = uint8(lTemp * 255)
	} else {
		if lTemp < 0.5 {
			var_2 = lTemp * ( 1 + sTemp)
		} else {
			var_2 = (lTemp + sTemp) - (sTemp * lTemp)
		}

		var_1 = 2 * lTemp - var_2
		r = uint8(255 * hue2rgb(var_1, var_2, hTemp+(1.0/3.0)))
		g = uint8(255 * hue2rgb(var_1, var_2, hTemp))
		b = uint8(255 * hue2rgb(var_1, var_2, hTemp - (1.0/3.0)))
	}
	return r, g, b
}

func hue2rgb(var_1, var_2, var_H float32) float32{
	v1, v2, vH := var_1, var_2, var_H
	if vH < 0 {
		vH += 1
	}
	if vH > 1 {
		vH -= 1
	}
	if 6 * vH < 1 {
		return (v1 + (v2 - v1) * 6 * vH)
	}

	if 2 *vH < 1 {
		return v2
	}

	if 3 * vH < 2 {
		return v1 + (v2 - v1) * (( 2.0/3.0) - vH ) * 6
	}

	return v1
}

func openImage(fileName string, isLocal bool) (*os.File, error){
	imagePath := getFullPath(fileName, isLocal)
	return os.Open(imagePath)
}

func min(x, y int) int{
	if x < y {
		return x
	}
	return y
}

func runImageWorkers(worker interface{}, bounds image.Rectangle) {
	switch worker.(type) {
	case func(int, int):
		start := 0
		chunkSize := bounds.Max.Y / numThreads
		for i := 0 ; i < numThreads ; i ++ {
			end := min(start + chunkSize, bounds.Max.Y)
			go worker.(func(int,int))(start, end)
			start = start + chunkSize
		}
	}
}

func HSLHistogramEqualizationConcurrent(fileName string, isLocal bool){
	img, err := openImage(fileName, isLocal)

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()

	// convert image from rgb to hsl, storing pixel values in arrays instead of temp Image
	imgSize := bounds.Max.X * bounds.Max.Y
	h_img := make([]float32, imgSize)
	s_img := make([]float32, imgSize)
	l_img := make([]uint8, imgSize)
	a_img := make([]uint8, imgSize)


	wg := sync.WaitGroup{}
	wg.Add(numThreads)

	conversionWorkers := func(start, end int) {
		defer wg.Done()
		for y := start; y < end; y ++ {
			for x := bounds.Min.X; x < bounds.Max.X ; x++ {
				r, g, b, a := decodedImg.At(x, y).RGBA()
				// A color's RGBA method returns values in the range [0, 65535].
				// Shifting by 8 reduces this to the range [0, 255].
				r, g, b, a = r>>8, g>>8, b>>8, a>>8

				h, s, l := rgb2hsl(uint8(r), uint8(g), uint8(b))
				h_img[pixelVal(x,y, bounds)] = h
				s_img[pixelVal(x,y, bounds)] = s
				l_img[pixelVal(x,y, bounds)] = l
				a_img[pixelVal(x,y, bounds)] = uint8(a) // save for later.
			}
		}
	}

	runImageWorkers(conversionWorkers, bounds)
	wg.Wait()


	// create a lightness histogram
	l_hist := [256]uint32{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			l_hist[l_img[pixelVal(x,y, bounds)]]++
		}
	}

	// get lightness look up table
	lLUT := getLookUpTable(l_hist, imgSize)

	//generate new contrast enhanced image with lightness Look up tables
	newImg := createNewImage(bounds)

	wg = sync.WaitGroup{}
	wg.Add(numThreads)

	imgWriteWorker := func(start, end int) {
		defer wg.Done()
		for y := start; y < end; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r,g,b := hsl2rgb(
					h_img[pixelVal(x,y,bounds)],
					s_img[pixelVal(x,y,bounds)],
					lLUT[l_img[pixelVal(x,y,bounds)]],  // new lightness value
				)


			newImg.Set(x,y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a_img[pixelVal(x,y, bounds)],
				})
			}
		}
	}

	runImageWorkers(imgWriteWorker, bounds)
	wg.Wait()

	f , imgError := writeImage(fileName, "enhanced-HSL-", isLocal)
	defer f.Close()
	if imgError != nil {
		panic(imgError)
	}

	if err != nil {
		panic(err)
	}

	defer f.Close()
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

func HSLHistogramEqualizationSerial(fileName string, isLocal bool){
	img, err := openImage(fileName, isLocal)

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()

	// convert image from rgb to hsl, storing pixel values in arrays instead of temp Image
	imgSize := bounds.Max.X * bounds.Max.Y
	h_img := make([]float32, imgSize)
	s_img := make([]float32, imgSize)
	l_img := make([]uint8, imgSize)
	a_img := make([]uint8, imgSize)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			r, g, b, a = r>>8, g>>8, b>>8, a>>8

			h, s, l := rgb2hsl(uint8(r), uint8(g), uint8(b))
			h_img[pixelVal(x,y, bounds)] = h
			s_img[pixelVal(x,y, bounds)] = s
			l_img[pixelVal(x,y, bounds)] = l
			a_img[pixelVal(x,y, bounds)] = uint8(a) // save for later.
		}
	}

	// create a lightness histogram
	l_hist := [256]uint32{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			l_hist[l_img[pixelVal(x,y, bounds)]]++
		}
	}

	// get lightness look up table
	lLUT := getLookUpTable(l_hist, imgSize)

	//generate new contrast enhanced image with lightness Look up tables
	newImg := createNewImage(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// convert back to rgb with hsl values
			r,g,b := hsl2rgb(
				h_img[pixelVal(x,y,bounds)],
				s_img[pixelVal(x,y,bounds)],
				lLUT[l_img[pixelVal(x,y,bounds)]],  // new lightness value
				)


			newImg.Set(x,y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a_img[pixelVal(x,y, bounds)],
			})
		}
	}

	f , imgError := writeImage(fileName, "enhanced-HSL-", isLocal)
	defer f.Close()
	if imgError != nil {
		panic(imgError)
	}

	if err != nil {
		panic(err)
	}

	defer f.Close()
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

func YUVHistogramEqualization(fileName string, isLocal bool){
	img, err := openImage(fileName, isLocal)

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()

	// convert image from rgb to yuv
	imgSize := bounds.Max.X * bounds.Max.Y

	// instead of creating a temp yuv image, lets keep track of the pixel values in arrays.
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

			// convert from rgb -> yuv / ycbcr
			Y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
			y_img[pixelVal(x,y, bounds)] = Y
			cb_img[pixelVal(x,y, bounds)] = cb
			cr_img[pixelVal(x,y, bounds)] = cr
			a_img[pixelVal(x,y, bounds)] = uint8(a)
		}
	}

	// We are actually only interested in Y channel (Luminance) for contrast enhancement.
	// 2 chrominance components may be ignored.
	y_hist := [256]uint32{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			y_hist[y_img[pixelVal(x,y, bounds)]]++
		}
	}

	yLUT := getLookUpTable(y_hist, imgSize)

	//generate new contrast enhanced image with y Look up tables
	newImg := createNewImage(bounds)
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

	f , imgError := writeImage(fileName, "enhanced-YUV-", isLocal)
	defer f.Close()
	if imgError != nil {
		panic(imgError)
	}

	//err = png.Encode(f, newImg)
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

func createNewImage(bounds image.Rectangle) *image.RGBA64{
	width, height := bounds.Max.X , bounds.Max.Y
	rect := image.Rect(0,0,width, height)
	newImg := image.NewRGBA64(rect)
	return newImg
}

// RGB contrast enhancement
func RGBHistogramEqualizationSerial(fileName string, isLocal bool){
	img, err := openImage(fileName, isLocal)

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()

	// generate histogram for each RGB channel.
	rHistogram, gHistogram, bHistogram := [256]uint32{}, [256]uint32{}, [256]uint32{}

	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	// https://golang.org/pkg/image/
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()

			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			// 2^16 / 2^8 -> 2^8
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			rHistogram[r]++
			gHistogram[g]++
			bHistogram[b]++
		}
	}

	// construct the Look Up Table for rgb values
	imgSize := bounds.Max.Y * bounds.Max.X
	rLUT := getLookUpTable(rHistogram, imgSize)
	gLUT := getLookUpTable(gHistogram, imgSize)
	bLUT := getLookUpTable(bHistogram, imgSize)

	//generate new contrast enhanced image with Look up tables
	newImg := createNewImage(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			newImg.Set(x,y, color.RGBA{
				R: rLUT[r],
				G: gLUT[g],
				B: bLUT[b],
				A: uint8(a),
			})
		}
	}

	f , imgError := writeImage(fileName, "enhanced-RGB-", isLocal)
	defer f.Close()
	if imgError != nil {
		panic(imgError)
	}

	//err = png.Encode(f, newImg)
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

// RGB contrast enhancement
func RGBHistogramEqualizationConcurrent(fileName string, isLocal bool){
	img, err := openImage(fileName, isLocal)

	if err != nil {
		panic(err)
	}
	defer img.Close()

	decodedImg, _, err := image.Decode(img)
	if err != nil {
		panic(err)
	}

	bounds := decodedImg.Bounds()

	// generate histogram for each RGB channel.
	rHistogram, gHistogram, bHistogram := [256]uint32{}, [256]uint32{}, [256]uint32{}

	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	// https://golang.org/pkg/image/
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()

			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			// 2^16 / 2^8 -> 2^8
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			rHistogram[r]++
			gHistogram[g]++
			bHistogram[b]++
		}
	}

	// construct the Look Up Table for rgb values
	imgSize := bounds.Max.Y * bounds.Max.X
	rLUT := getLookUpTable(rHistogram, imgSize)
	gLUT := getLookUpTable(gHistogram, imgSize)
	bLUT := getLookUpTable(bHistogram, imgSize)

	//generate new contrast enhanced image with Look up tables
	newImg := createNewImage(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := decodedImg.At(x, y).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			newImg.Set(x,y, color.RGBA{
				R: rLUT[r],
				G: gLUT[g],
				B: bLUT[b],
				A: uint8(a),
			})
		}
	}

	f , imgError := writeImage(fileName, "enhanced-RGB-", isLocal)
	defer f.Close()
	if imgError != nil {
		panic(imgError)
	}

	//err = png.Encode(f, newImg)
	err = jpeg.Encode(f, newImg, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

func writeImage(fileName string,  prefix string, isLocal bool,) (*os.File, error) {
	var f *os.File
	var err error
	if isLocal {
		f, err = os.Create(prefix + fileName )
	} else {
		f, err = os.Create(VagrantImageDir + prefix + fileName )
	}
	return f, err
}

// create conversion values by constructing look up table with CDF
func getLookUpTable(histogram [256]uint32, imageSize int) [256]uint8 {
	// construct look up table by caluclating CDF
	sum := uint64(0)
	for _, val := range histogram {
		sum += uint64(val)
	}

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
	return lut
}
