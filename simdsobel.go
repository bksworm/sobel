package sobel

/*
#cgo pkg-config: libsimd
#include <stdlib.h>
#include "Simd/SimdLib.h"
extern int  sobelSimdGray8(uint8_t *src, size_t width, size_t height, uint8_t *dst)  ;
*/
import "C"

import (
	"image"
	"math"
	"reflect"
	"unsafe"
)

//BenchmarkIT/Benchmark_FilterGraySimd-2         	    1032	   1499811 ns/op malloc
//
//BenchmarkIT/Benchmark_FilterGraySimd-2         	     501	   2000184 ns/op
// dstX := make([]uint16, dstSize) dstXC := (*C.uint8_t)(unsafe.Pointer(&dstX[0]))
//
//BenchmarkIT/Benchmark_FilterGraySimd-2         	     276	   4853208 ns/op
//add flited filling
func FilterGraySimd(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min
	filtered = image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))
	imSize := len(grayImg.Pix)

	// void SimdSobelDyAbs (
	//    [in]	src	- a pointer to pixels data of the input image.
	//    [in]	srcStride	- a row size of the input image.
	//    [in]	width	- an image width.
	//    [in]	height	- an image height.
	//    [out]	dst	- a pointer to pixels data of the output image.
	//    [in]	dstStride	- a row size of the output image (in bytes).
	// )
	//All images must have the same width and height. Input image
	// must has 8-bit gray format, output image must has 16-bit integer format.
	src := (*C.uint8_t)(unsafe.Pointer(&grayImg.Pix[0]))
	srcStride := C.size_t(grayImg.Stride)
	width := C.size_t(max.X)
	height := C.size_t(max.Y)
	dstStride := srcStride * 2
	dstSize := C.size_t(imSize * 2)

	dstXC := (*C.uint8_t)(unsafe.Pointer(C.malloc(dstSize)))
	defer C.free(unsafe.Pointer(dstXC))
	C.SimdSobelDxAbs(src, srcStride, width, height, dstXC, dstStride)

	dstYC := (*C.uint8_t)(unsafe.Pointer(C.malloc(dstSize)))
	defer C.free(unsafe.Pointer(dstYC))

	C.SimdSobelDyAbs(src, srcStride, width, height, dstYC, dstStride)

	dstX := nonCopyGoUint16(uintptr(unsafe.Pointer(dstXC)), imSize)
	dstY := nonCopyGoUint16(uintptr(unsafe.Pointer(dstYC)), imSize)

	var fX, fY uint
	var fS float64
	var pix uint8
	if len(dstX) == imSize && len(dstY) == imSize && len(filtered.Pix) == imSize {
		for i := 0; i < imSize; i++ {
			fX = uint(dstX[i])
			fY = uint(dstY[i])
			fS = math.Sqrt(float64(fX*fX + fY*fY))
			//clipping
			if fS > 255.0 {
				pix = 255
			} else {
				pix = uint8(fS)
			}
			filtered.Pix[i] = pix
		}
	}
	return filtered
}

func nonCopyGoUint16(ptr uintptr, length int) []uint16 {
	var slice []uint16
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = length
	header.Len = length
	header.Data = ptr
	return slice
}

func FilterGraySimdC(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min
	filtered = image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))
	src := (*C.uint8_t)(unsafe.Pointer(&grayImg.Pix[0]))
	dst := (*C.uint8_t)(unsafe.Pointer(&filtered.Pix[0]))
	width := C.size_t(max.X)
	height := C.size_t(max.Y)
	C.sobelSimdGray8(src, width, height, dst)

	return filtered
}
