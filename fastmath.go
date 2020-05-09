package sobel

import "math"

// ISqrt returns floor(sqrt(n)). Typical run time is few hundreds of ns.
//https://gitlab.com/cznic/mathutil/-/blob/master/mathutil.go
func ISqrt(n uint32) (x uint32) {
	if n == 0 {
		return
	}

	if n >= math.MaxUint16*math.MaxUint16 {
		return math.MaxUint16
	}
	var px, nx uint32
	for x = n; ; px, x = x, nx {
		nx = (x + n/x) / 2
		if nx == x || nx == px {
			break
		}
	}
	return
}

//https://www.geeksforgeeks.org/square-root-of-an-integer/
//Time Complexity: O(Log x)
func FloorSqrt(x uint32) (ans uint32) {
	// Base Cases
	if x == 0 || x == 1 {
		return x
	}

	// Do Binary Search for floor(sqrt(x))
	var (
		start uint32 = 1
		mid   uint32
	)
	end := x

	for start <= end {
		mid = (start + end) / 2

		// If x is a perfect square
		if mid*mid == x {
			return mid
		}

		// Since we need floor, we update answer when mid*mid is
		// smaller than x, and move closer to sqrt(x)
		if mid*mid < x {

			start = mid + 1
			ans = mid
		} else { // If mid*mid is greater than x
			end = mid - 1
		}
	}
	return ans
}

//Note: The Binary Search can be further optimized to start with ‘start’ = 0 and ‘end’ = x/2.
//Floor of square root of x cannot be more than x/2 when x > 1.
//Benchmark_FloorSqrt-8       25.7 ns/op
//Benchmark_FloorSqrtFast-8   20.5 ns/op

func FloorSqrtFast(x uint32) (ans uint32) {
	// Base Cases
	if x == 0 || x == 1 {
		return x
	}

	// Do Binary Search for floor(sqrt(x))
	var (
		start uint32 = 0
		mid   uint32
	)
	end := x / 2

	for start <= end {
		mid = (start + end) / 2

		// If x is a perfect square
		if mid*mid == x {
			return mid
		}

		// Since we need floor, we update answer when mid*mid is
		// smaller than x, and move closer to sqrt(x)
		if mid*mid < x {

			start = mid + 1
			ans = mid
		} else { // If mid*mid is greater than x
			end = mid - 1
		}
	}
	return ans
}

// Abs returns the absolute value of the given int.
func Abs(x int) uint32 {
	if x < 0 {
		return uint32(-x)
	} else {
		return uint32(x)
	}
}
