package main
import (
   "fmt"
   "time"
   "gobot.io/x/gobot"
   "gobot.io/x/gobot/platforms/dji/tello"
)
//cannot have 2 mains in same folder
//func main() {
func swap() {
   fmt.Println("Start")
   drone := tello.NewDriver("8888")
   work := func() {
      drone.TakeOff()
      gobot.After(3*time.Second, func() {
         drone.Left(50)
         time.Sleep(time.Second*3)
         drone.Right(50)
         time.Sleep(time.Second*3)
         drone.Land()
      })
   }
   robot := gobot.NewRobot("tello",
      []gobot.Connection{},
      []gobot.Device{drone},
      work,
   )
   robot.Start()
}
