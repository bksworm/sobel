package sobel

import (
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"testing"

	"github.com/bksworm/sobel/testsuite"
)

type SobelTS struct {
	tst   *testing.T
	bench *testing.B
	img   *image.Gray
}

const fileName = "../img/test.png"

func decodePng(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// SetUpSuite is called once before the very first test in suite runs
func (s *SobelTS) SetUpSuite() {

	imgSrc, err := decodePng(fileName)
	if err != nil {
		s.tst.Fatalf("%s: %v", fileName, err)
	}
	s.img = ToGrayscale(imgSrc)
}

// TearDownSuite is called once after thevery last test in suite runs
func (s *SobelTS) TearDownSuite() {
}

// SetUp is called before each test method
func (s *SobelTS) SetUp() {
}

// TearDown is called after each test method
func (s *SobelTS) TearDown() {
}

// Hook up  into the "go test" runner.
func TestIt(t *testing.T) {
	s := &SobelTS{tst: t}
	testsuite.Run(t, s)
}

func BenchmarkIT(b *testing.B) {
	s := &SobelTS{bench: b}
	testsuite.Bench(b, s)
}

//THERE ARE TEST

func (s *SobelTS) Benchmark_applySobelFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		applySobelFilter(s.img, 7, 7)
	}
}

func (s *SobelTS) Benchmark_applySobelFilterFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		applySobelFilterFast(s.img, 7, 7)
	}
}

func (s *SobelTS) Benchmark_FilterGray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterGray(s.img, SobelFast)
	}
}

func (s *SobelTS) Benchmark_FilterGrayFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterGrayFast(s.img, SobelFast)
	}
}

func (s *SobelTS) Benchmark_FilterGrayMath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterGrayMath(s.img)
	}
}

func (s *SobelTS) Benchmark_FilterGraySimd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterGraySimd(s.img)
	}
}

func (s *SobelTS) Benchmark_FilterGraySimdC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterGraySimdC(s.img)
	}
}

const sqrtFrom = 4356789
const runsNumber = 10000

func (s *SobelTS) Test_ISqrt(t *testing.T) {
	log.Printf("sqrt( %d ) = %d ", sqrtFrom, ISqrt(sqrtFrom))
}

func (s *SobelTS) Test_FloorSqrt(t *testing.T) {
	log.Printf("sqrt( %d ) = %d ", sqrtFrom, FloorSqrt(sqrtFrom))
}

func (s *SobelTS) Test_FloorSqrtFast(t *testing.T) {
	log.Printf("sqrt( %d ) = %d ", sqrtFrom, FloorSqrtFast(sqrtFrom))
}

func (s *SobelTS) Test_FilterGraySimd(t *testing.T) {
	for i := 0; i < 100; i++ {
		FilterGraySimd(s.img)
	} //
}

func (s *SobelTS) Benchmark_SqrtI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ISqrt(sqrtFrom)
	}
}

func (s *SobelTS) Benchmark_SqrtFloor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FloorSqrt(sqrtFrom)
	}
}

func (s *SobelTS) Benchmark_SqrtFloorFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FloorSqrtFast(sqrtFrom)
	}
}

func (s *SobelTS) Benchmark_SqrtFloorC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FloorSqrtC(sqrtFrom)
	}
}

func (s *SobelTS) Benchmark_SqrtMath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		math.Sqrt(float64(sqrtFrom))
	}
}

func (s *SobelTS) Benchmark_Abs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abs(-3)
	}
}
