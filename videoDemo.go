package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
)

const (
	frameSize = 960 * 720 * 3
)

func main() {
	drone := tello.NewDriver("8890")
	//window := opencv.NewWindowDriver()
	window := gocv.NewWindow("Demo2")
	classifier := gocv.NewCascadeClassifier()
	classifier.Load("haarcascade_frontalface_default.xml")
	defer classifier.Close()
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()
	work := func() {

		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}
		//count:=0
		go func() {

		}()

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetExposure(1)
			drone.SetVideoEncoderRate(4)

			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start(false)
	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if img.Empty() {
			continue
		}
		imageRectangles := classifier.DetectMultiScale(img)

		for _, rect := range imageRectangles {
			log.Println("found a face,", rect)
			gocv.Rectangle(&img, rect, colornames.Cadetblue, 3)
		}
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
		//if count < 1000{
		//
		//}
		//window.WaitKey(1)
		//count +=1
	}
}