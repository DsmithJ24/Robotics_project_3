package main

import (
	"fmt" // Formatted I/O
	"io" //  It provides basic interfaces to I/O primitives
	"os/exec" // To run the external commands.
	"strconv" // Package strconv implements conversions to and from string
	"time" //For time related operation
	"log"
    "image"
    "image/color"
    "os"
    "math"

	"gobot.io/x/gobot" // Gobot Framework.
	"gobot.io/x/gobot/platforms/dji/tello" // DJI Tello package.
	"gocv.io/x/gocv" // GoCV package to access the OpenCV library.
	"golang.org/x/image/colornames"
)

var tracking = false
var detectSize = false
var distTolerance = 0.05 * dist(0, 0, frameX, frameY)

// Frame size constant.
const (
	frameX    = 720
	frameY    = 960
	frameSize = frameX * frameY * 3
)

func main() {
	// Driver: Tello Driver
	drone := tello.NewDriver("8890")

    //window := opencv.NewWindowDriver()
	// OpenCV window to watch the live video stream from Tello.
	window := gocv.NewWindow("Tello")

	//classifier stuff
	classifier := gocv.NewCascadeClassifier()
	classifier.Load("haarcascade_frontalface_default.xml")
	defer classifier.Close()

//	classifier.Load("eyelook.xml")
	//blue := color.RGBA{0, 0, 255, 0}
//	defer classifier.Close()

	//FFMPEG command to convert the raw video from the drone.
	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameY)+"x"+strconv.Itoa(frameX), "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()



	work := func() {
		//Starting FFMPEG.
		fmt.Println("starting ffmpeg")
		if err := ffmpeg.Start(); err != nil {
			fmt.Println("there is an error")
			fmt.Println(err)
			return
		}

        // need this
        go func() {

        }()

		// Event: Listening the Tello connect event to start the video streaming.
		go drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected to Tello.")
			drone.StartVideo()
			drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			drone.SetExposure(0)

			//For continued streaming of video.
			gobot.Every(500*time.Millisecond, func() {
				drone.StartVideo()

			})
		})

		//Event: Piping the video data into the FFMPEG function.
		go drone.On(tello.VideoFrameEvent, func(data interface{}) {
			fmt.Println("receiving data")
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

		//TakeOff the Drone.
		gobot.After(5*time.Second, func() {
			go drone.TakeOff()
			fmt.Println("Tello Taking Off...")
		})

		//Land the Drone.
		gobot.After(30*time.Second, func() {
			go drone.Land()
			fmt.Println("Tello Landing...")
		})


	}

	//Robot: Tello Drone
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	// calling Start(false) lets the Start routine return immediately without an additional blocking goroutine
	robot.Start(false)

	if len(os.Args) < 3 {
		//fmt.Println("How to run:\ngo run facetracking.go [protofile] [modelfile]")
		return
	}

    fmt.Println("test")
	proto := os.Args[1]
	model := os.Args[2]

	net := gocv.ReadNetFromCaffe(proto, model)
	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", proto, model)
		return
	}
	defer net.Close()

	//green := color.RGBA{0, 255, 0, 0}

	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", proto, model)
		return
	}

	green := color.RGBA{0, 255, 0, 0}
	defer net.Close()
	//treacking stuff here
	refDistance := float64(0)
    detected := false
    left := float32(0)
    top := float32(0)
    right := float32(0)
    bottom := float32(0)

	// now handle video frames from ffmpeg stream in main thread, to be macOs friendly
	for {
		buf := make([]byte, frameSize)
//		fmt.Println("handle vid frames")
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
//			fmt.Println("error, jumping up")
			//EOF is error being printed, fixed as of 11/16
			fmt.Println(err)
			continue
		}

		img, _ := gocv.NewMatFromBytes(frameX, frameY, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
		    fmt.Println("Empty image")
			continue
		}

//face tracking
		W := float32(img.Cols())
		H := float32(img.Rows())
		blob := gocv.BlobFromImage(img, 1.0, image.Pt(128, 96), gocv.NewScalar(104.0, 177.0, 123.0, 0), false, false)
		defer blob.Close()

		net.SetInput(blob, "data")

		detBlob := net.Forward("detection_out")
		defer detBlob.Close()

		detections := gocv.GetBlobChannel(detBlob, 0, 0)
		defer detections.Close()

		for r := 0; r < detections.Rows(); r++ {
			confidence := detections.GetFloatAt(r, 2)
			if confidence < 0.5 {
				continue
			}

			left = detections.GetFloatAt(r, 3) * W
			top = detections.GetFloatAt(r, 4) * H
			right = detections.GetFloatAt(r, 5) * W
			bottom = detections.GetFloatAt(r, 6) * H

			left = min(max(0, left), W-1)
			right = min(max(0, right), W-1)
			bottom = min(max(0, bottom), H-1)
			top = min(max(0, top), H-1)

			rect := image.Rect(int(left), int(top), int(right), int(bottom))
			gocv.Rectangle(&img, rect, green, 3)
			detected = true
		}

//face detect
		faceDetect := gocv.NewMat()
		gocv.Resize( img, &faceDetect, image.Pt( 90, 120 ), 0, 0, gocv.InterpolationNearestNeighbor)

		//bad, too grainy
        //gocv.Resize( img, &faceDetect, image.Pt( 70, 100 ), 0, 0, gocv.InterpolationNearestNeighbor)

		//detect a face
		imageRectangles := classifier.DetectMultiScale( faceDetect )

		for _, rect := range imageRectangles {
			log.Println("found a face,", rect)
			gocv.Rectangle(&faceDetect, rect, colornames.Cadetblue, 3)
		}

/*
		//for eyes
		for _, rect := range imageRectangles {
        			gocv.Rectangle(&img, rect, blue, 3)

        			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
        			pt := image.Pt(rect.Min.X+(rect.Min.X/2)-(size.X/2), rect.Min.Y-2)
        			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
        			return
        		}
*/
		window.IMShow(faceDetect)
		if window.WaitKey(1) >= 0 {
			break
		}

		//movement for face tracking
		if !tracking || !detected {
    		continue
    	}

    	if detectSize {
    		detectSize = false
    		refDistance = dist(left, top, right, bottom)
    	}

    	distance := dist(left, top, right, bottom)

    	if right < W/2 {
    		drone.CounterClockwise(50)
    	} else if left > W/2 {
    		drone.Clockwise(50)
    	} else {
    		drone.Clockwise(0)
    	}

    	if top < H/10 {
    		drone.Up(25)
    	} else if bottom > H-H/10 {
    		drone.Down(25)
    	} else {
    		drone.Up(0)
    	}

    	if distance < refDistance-distTolerance {
    		drone.Forward(20)
    	} else if distance > refDistance+distTolerance {
    		drone.Backward(20)
    	} else {
    		drone.Forward(0)
    	}
	}
}

func dist(x1, y1, x2, y2 float32) float64 {
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}