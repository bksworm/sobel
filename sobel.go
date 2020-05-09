package sobel

import (
	"image"
	"image/color"
	"math"
)

var (
	sobelX = [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	sobelY = [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	sobelXF = [9]float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	sobelYF = [9]float64{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	sobelXL = [9]int{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	sobelYL = [9]int{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	sharaX = [3][3]int{
		{-3, 0, 3},
		{-10, 0, 10},
		{-3, 0, 3},
	}
	sharaY = [3][3]int{
		{3, 10, 3},
		{0, 0, 0},
		{-3, -10, -3},
	}

	laplasianX = [3][3]int{
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}
	laplasianY = [3][3]int{
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}

	sharpenX = [3][3]int{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
	sharpenY = [3][3]int{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
)

type FilterType int

const kernelSize = 3

const (
	Sobel FilterType = iota
	SobelFast
	Laplasian
	Shara
	Sharpen
)

type filterFunc func(img *image.Gray, x int, y int) (uint32, uint32)

func Filter(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGrayFast(grayImg, flt)
}

func FilterMath(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGrayMath(grayImg)
}

func FilterSimd(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGraySimd(grayImg)
}

func getFilterFunc(flt FilterType) (res filterFunc) {
	switch flt {
	case Sobel:
		res = applySobelFilter
		break
	case SobelFast:
		res = applySobelFilterFast
		break
	case Laplasian:
		res = applyLaplasianFilter
		break
	case Shara:
		res = applySharaFilter
		break
	case Sharpen:
		res = applySharpenFilter
	}
	return res
}

//for better optimization in case of input gray image
func FilterGray(grayImg *image.Gray, flt FilterType) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min
	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide 1 pixel "border"
	height := max.Y - 1
	applay := getFilterFunc(flt)
	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applay(grayImg, x, y)
			v := ISqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			pixel := color.Gray{Y: uint8(v)}
			filtered.SetGray(x, y, pixel)
		}
	}
	return filtered
}

//Benchmark_FilterGray	   27336411 ns/op
//Benchmark_FilterGrayFast 19521755 ns/op
//for better optimization in case of input gray image
func FilterGrayFast(grayImg *image.Gray, flt FilterType) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	var v uint32
	applay := getFilterFunc(flt)

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applay(grayImg, x, y)
			v = FloorSqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			filtered.SetGray(x, y, color.Gray{Y: uint8(v)})
		}
	}

	return filtered
}

func FilterGrayMath(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applySobelFilterMath(grayImg, x, y)
			fS := (fX * fX) + (fY * fY)
			//filtered.SetGray(x, y, color.Gray{Y: uint8(v)})
			filtered.Pix[filtered.PixOffset(x-1, y-1)] = uint8(math.Sqrt(fS))
		}
	}

	return filtered
}

func applyLaplasianFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += laplasianX[i][j] * pixel
			fY += laplasianY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySobelFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY, pixel int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel = int(img.GrayAt(curX, curY).Y)
			fX += sobelX[i][j] * pixel
			fY += sobelY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

//Benchmark_applySobelFilter	   	33.8 ns/op
//Benchmark_applySobelFilterFast   	19.0 ns/op
func applySobelFilterFast(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		//index = i * kernelSize
		for j := 0; j < kernelSize; j++ {
			//it is unsafe but faster on 10% or so
			pixel = int(img.Pix[img.PixOffset(curX, curY)])
			fX += sobelXL[index] * pixel
			fY += sobelYL[index] * pixel
			curX++
			index++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySobelFilterMath(img *image.Gray, x int, y int) (float64, float64) {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		//index = i * kernelSize
		for j := 0; j < kernelSize; j++ {
			//it is unsafe but faster on 10% or so
			pixel = int(img.Pix[img.PixOffset(curX, curY)])
			fX += sobelXL[index] * pixel
			fY += sobelYL[index] * pixel
			curX++
			index++
		}
		curX = x - 1
		curY++
	}
	return math.Abs(float64(fX)), math.Abs(float64(fY))
}

func applySharaFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += sharaX[i][j] * pixel
			fY += sharaY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySharpenFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += sharpenX[i][j] * pixel
			fY += sharpenY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}
