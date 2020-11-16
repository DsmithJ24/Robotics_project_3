package main
import (
   "fmt"
   //"os/exec"
   "time"
   "gobot.io/x/gobot"
   "gobot.io/x/gobot/platforms/dji/tello"

   "gocv.io/x/gocv"
//   "golang.org/x/image/colornames"
)
//from intelliJ

const (
	frameSize = 960 * 720 * 3
)

func main() {
	drone := tello.NewDriver("8890")

	//this is the take off function
	takeOff := func() {
		fmt.Println("Taking off")
		drone.TakeOff()

    	gobot.After(10*time.Second, func() {
        	drone.Land()
        	fmt.Println("Landed")
		})
	}

/*
//this is from the new ffmpeg demo code:
    window := gocv.NewWindow("Demo2")
    classifier := gocv.NewCascadeClassifier()
   	classifier.Load("haarcascade_frontalface_default.xml")
    defer classifier.Close()
//    ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
    	"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
    ffmpegIn, _ := ffmpeg.StdinPipe()
    ffmpegOut, _ := ffmpeg.StdoutPipe()

    takeOff := func() {
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

   */
   robot := gobot.NewRobot("tello",
      []gobot.Connection{},
      []gobot.Device{drone},
      takeOff,
   )
   robot.Start()
}