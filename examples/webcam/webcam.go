// Example program that uses blakjack/webcam library
// for working with V4L2 devices.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/pprof" //add profiling
	"net/textproto"
	"sort"
	"strconv"
	"time"

	"image/jpeg"

	"github.com/bksworm/sobel"
	"github.com/blackjack/webcam"
)

const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_YUYV = 0x56595559
	V4L2_PIX_FMT_MJPG = 0x47504a4d
)

type FrameSizes []webcam.FrameSize

func (slice FrameSizes) Len() int {
	return len(slice)
}

//For sorting purposes
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

//For sorting purposes
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
	V4L2_PIX_FMT_MJPG: true,
}

func main() {
	dev := flag.String("d", "/dev/video0", "video device to use")
	fmtstr := flag.String("f", "", "video format to use, default first supported")
	// szstr := flag.String("s", "", "frame size to use, default largest one")
	// single := flag.Bool("m", false, "single image http mode, default mjpeg video")
	addr := flag.String("l", ":8080", "addr to listien")
	fps := flag.Bool("p", false, "print fps info")
	flag.Parse()

	cam, err := webcam.Open(*dev)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer cam.Close()

	// select pixel format
	format_desc := cam.GetSupportedFormats()

	fmt.Println("Available formats:")
	for _, s := range format_desc {
		fmt.Println(s)
	}

	//input format selection
	var format webcam.PixelFormat
FMT:
	for f, s := range format_desc {
		if *fmtstr == "" {
			if supportedFormats[f] {
				format = f
				break FMT
			}

		} else if *fmtstr == s {
			if !supportedFormats[f] {
				log.Println(format_desc[f], "format is not supported, exiting")
				return
			}
			format = f
			break
		}
	}

	if format == 0 {
		log.Println("No format found, exiting")
		return
	}

	// select frame size
	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)
	fmt.Println("Supported frame sizes for format", format_desc[format])
	for _, f := range frames {
		fmt.Println(f.GetString())
	}

	f, w, h, err := cam.SetImageFormat(format, 640, 480) //uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		log.Println("SetImageFormat return error", err)
		return
	}

	fmt.Printf("Resulting image format: %s %dx%d\n", format_desc[f], w, h)

	// start streaming
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return
	}

	var (
		li   chan *bytes.Buffer = make(chan *bytes.Buffer)
		fi   chan []byte        = make(chan []byte)
		back chan struct{}      = make(chan struct{})
	)

	go encodeToImage(cam, back, fi, li, w, h, f)
	go httpVideo(*addr, li)

	timeout := uint32(5) //5 seconds
	start := time.Now()
	var fr time.Duration

	for {
		err = cam.WaitForFrame(timeout)
		if err != nil {
			log.Println(err)
			return
		}

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			log.Println(err)
			continue
		default:
			log.Println(err)
			return
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Println(err)
			return
		}

		if len(frame) != 0 {
			// print framerate info every 10 seconds
			fr++
			if *fps {
				d := time.Since(start)
				if d > time.Second*10 {
					fmt.Println(float64(fr)/(float64(d)/float64(time.Second)), "fps")
					start = time.Now()
					fr = 0
				}
			}
			//send frame pointer
			select {
			case fi <- frame:
				<-back
			default:
			}
		}
	}
}

func encodeToImage(wc *webcam.Webcam, back chan struct{}, fi chan []byte,
	li chan *bytes.Buffer, w, h uint32, format webcam.PixelFormat) {

	var (
		frame []byte
		buf   *bytes.Buffer
	)

	grayImg := image.NewGray(image.Rect(0, 0, int(w), int(h)))
	for {
		bframe := <-fi
		//make sure that we have enough memory
		if len(frame) < len(bframe) {
			frame = make([]byte, len(bframe))
		}

		copy(frame, bframe) // copy frame
		back <- struct{}{}  //we are redy to receive a new frame

		switch format {
		case V4L2_PIX_FMT_YUYV:
			for i := range grayImg.Pix {
				i2 := i * 2
				grayImg.Pix[i] = frame[i2] //copy only Y set
			}
			//select edges
			edges := sobel.FilterGray(grayImg, sobel.Sobel)
			//convert to jpeg
			buf = &bytes.Buffer{}
			if err := jpeg.Encode(buf, edges, nil); err != nil {
				log.Fatal(err)
				return
			}

		case V4L2_PIX_FMT_MJPG:
			buf = bytes.NewBuffer(frame)

		default:
			log.Fatal("invalid format ?")
		}

		const N = 50
		// broadcast image up to N ready clients
		nn := 0
	FOR:
		for ; nn < N; nn++ {
			select {
			case li <- buf:
			default:
				break FOR
			}
		}

		if nn == 0 {
			li <- buf
		}

	}
}

func httpVideo(addr string, li chan *bytes.Buffer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("connect from", r.RemoteAddr, r.URL)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		//remove stale image
		<-li
		const boundary = `frame`
		w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
		multipartWriter := multipart.NewWriter(w)
		multipartWriter.SetBoundary(boundary)
		for {
			img := <-li
			image := img.Bytes()
			iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
				"Content-type":   []string{"image/jpeg"},
				"Content-length": []string{strconv.Itoa(len(image))},
			})
			if err != nil {
				log.Println(err)
				return
			}
			_, err = iw.Write(image)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
