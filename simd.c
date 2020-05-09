#include <stdlib.h>
#include <math.h>
#include "Simd/SimdLib.h"

#define TRUE 1
#define FALSE 0

int  sobelSimdGray8(uint8_t *src, size_t width, size_t height, uint8_t *dst) {
    size_t  imgSize = width * height ;
	size_t srcStride = width ;
    size_t dstStride = srcStride * 2 ;
    size_t dstSize = imgSize * 2 ;

    uint8_t * dstXC = malloc(dstSize) ;
    if (NULL == dstXC ) {
        return FALSE ;
    }
    uint8_t *dstX = dstXC ;

    uint8_t * dstYC = malloc(dstSize) ;
    if (NULL == dstYC ) {
        free(dstXC) ;
        return FALSE ;
    }
    uint8_t *dstY = dstYC ;

    SimdSobelDxAbs(src, srcStride, width, height, dstX, dstStride) ;
    SimdSobelDyAbs(src, srcStride, width, height, dstY, dstStride) ;

    uint32_t fX, fY ;
    double  fS ;
    uint8_t pix ;
    for (int i = 0; i < imgSize; i++ ) {
        fX = (uint32_t)*dstX ;
        fY = (uint32_t)*dstY ;
        fS = sqrt((double)(fX*fX + fY*fY));
        //clipping
        if (fS > 255.0) {
            pix = 255 ;
        } else {
            pix = (uint8_t)fS ;
        }
        *dst = pix ;
        dst++ ;
        dstX ++ ;
        dstY ++ ;
    }

    free(dstXC) ;
    free(dstYC) ;
    return TRUE ;
}