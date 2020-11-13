package main

import (
	"fmt" // Formatted I/O
	"io" //  It provides basic interfaces to I/O primitives
	"os/exec" // To run the external commands.
	"strconv" // Package strconv implements conversions to and from string
	"time" //For time related operation

	"gobot.io/x/gobot" // Gobot Framework.
	"gobot.io/x/gobot/platforms/dji/tello" // DJI Tello package.
	"gocv.io/x/gocv" // GoCV package to access the OpenCV library.
)

// Frame size constant.
const (
	frameX    = 960
	frameY    = 720
	frameSize = frameX * frameY * 3
)

func main() {
	// Driver: Tello Driver
	drone := tello.NewDriver("8890")

	// OpenCV window to watch the live video stream from Tello.
	window := gocv.NewWindow("Tello")

	//FFMPEG command to convert the raw video from the drone.
	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")
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
		// Event: Listening the Tello connect event to start the video streaming.
		//does not connect
		fmt.Println("connecting...")
		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected to Tello.")
			drone.StartVideo()
			drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			drone.SetExposure(0)

			//For continued streaming of video.
			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		//Event: Piping the video data into the FFMPEG function.
		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			//fmt.Println("receiving data")
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
/*
		//TakeOff the Drone.
		gobot.After(5*time.Second, func() {
			drone.TakeOff()
			fmt.Println("Tello Taking Off...")
		})

		//Land the Drone.
		gobot.After(15*time.Second, func() {
			drone.Land()
			fmt.Println("Tello Landing...")
		})
*/
	}
	//Robot: Tello Drone
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	// calling Start(false) lets the Start routine return immediately without an additional blocking goroutine
	robot.Start(false)

	// now handle video frames from ffmpeg stream in main thread, to be macOs friendly
	for {
		buf := make([]byte, frameSize)
		fmt.Println("handle vid frames")
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}
		//no image, it is empty
		img, _ := gocv.NewMatFromBytes(frameY, frameX, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
		    fmt.Println("Empty image")
			continue
		}
		//never gets here due to above
        fmt.Println("Show vid")
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

/*
func handleConnected(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		fmt.Println("Drone connected.")

		connected = true

		drone.SetVideoEncoderRate(2)
		gobot.Every(100*time.Millisecond, func() {
			drone.StartVideo()
		})

		gobot.Every(time.Duration((1000.0/tickRate))*time.Millisecond, func() {
			if connected {
				tick(drone)
			}
		})
	}
}
*/