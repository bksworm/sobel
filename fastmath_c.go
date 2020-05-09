package sobel

/*
#include <stdint.h>

// Returns floor of square root of x
uint32_t floorSqrt(uint32_t x)
{
    // Base cases
    if (x == 0 || x == 1)
       return x;

    // Do Binary Search for floor(sqrt(x))
    uint32_t start = 0, end = x/2, ans;
    while (start <= end)
    {
        uint32_t mid = (start + end) / 2;

        // If x is a perfect square
        if (mid*mid == x)
            return mid;

        // Since we need floor, we update answer when mid*mid is
        // smaller than x, and move closer to sqrt(x)
        if (mid*mid < x)
        {
            start = mid + 1;
            ans = mid;
        }
        else // If mid*mid is greater than x
            end = mid-1;
    }
    return ans;
}
*/
import "C"

func FloorSqrtC(x uint32) uint32 {
	return uint32(C.floorSqrt(C.uint32_t(x)))
}
