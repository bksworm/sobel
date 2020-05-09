package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/bksworm/sobel"
)

const (
	outputFileMode = 644
)

func main() {
	inFile := flag.String("f", "", "input file")
	outFile := flag.String("o", "sobel", "ouput file")
	cpuprofile := flag.String("p", "", "write cpu profile to file")

	flag.Parse()
	if *inFile == "" {
		fmt.Printf("error: must provide input image")
		fmt.Println("usage: ")
		flag.PrintDefaults()
		os.Exit(0)
	}

	contents, err := os.Open(*inFile)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	defer contents.Close()

	img, ftype, err := image.Decode(contents)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		log.Fatalf("filetype %s not supported\n", ftype)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var edged = sobel.FilterSimd(img, sobel.Sobel)

	var ext string

	if ftype != "png" && ftype != "jpeg" {
		log.Fatalf("can't encode file type")
		os.Exit(0)
	} else if ftype == "png" {
		if !strings.Contains(*outFile, ".png") {
			ext = ".png"
		} else {
			ext = ""
		}
		out, err := os.Create(*outFile + ext)
		handleError(err)
		defer out.Close()
		err = png.Encode(out, edged)
		handleError(err)
	} else if ftype == "jpeg" {
		if !strings.Contains(*outFile, ".jpg") {
			ext = ".jpg"
		} else {
			ext = ""
		}
		out, err := os.Create(*outFile + ext)
		handleError(err)
		defer out.Close()
		err = jpeg.Encode(out, edged, nil)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
