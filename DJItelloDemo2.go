package main
import (
   "fmt"
   "os/exec"
   "time"
   "gobot.io/x/gobot"
   "gobot.io/x/gobot/platforms/dji/tello"
)
//ask why it does not do anything
func swap() {
//func main() {
   //this driver # is for video and flight
   drone := tello.NewDriver("8890")
   work := func() {
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
   }
   robot := gobot.NewRobot("tello",
      []gobot.Connection{},
      []gobot.Device{drone},
      work,
   )
   robot.Start()
}
