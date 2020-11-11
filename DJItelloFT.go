package main
import (
   "fmt"
   "os/exec"
   "time"
   "gobot.io/x/gobot"
   "gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8890")
	work := func() {
		fmt.Println("Taking off")
		drone.TakeOff()

		/*
		//mplayer sux
		mplayer := exec.Command("mplayer", "-fps", "25", "-")
    	mplayerIn, _ := mplayer.StdinPipe()
    	if err := mplayer.Start(); err != nil {
        	fmt.Println(err)
        	return
	  	}
	  	drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(4)
			gobot.Every(100*time.Millisecond, func() {
		   		drone.StartVideo()
			})
	 	})
		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
		//run multiple times?
			fmt.Println("Recieving Data")
			if _, err := mplayerIn.Write(pkt); err != nil {
		   		fmt.Println(err)
			}
	 	})
		 */
		
    	gobot.After(10*time.Second, func() {
        	drone.Land()
        	fmt.Println("Landed")
		})
		
   }
   robot := gobot.NewRobot("tello",
      []gobot.Connection{},
      []gobot.Device{drone},
      work,
   )
   robot.Start()
}