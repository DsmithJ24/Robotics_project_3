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
	drone.On(tello.FlightDataEvent, func(data interface{}) {
	   // TODO: protect flight data from race condition
	   flightData := data.(*tello.FlightData)
	   fmt.Println("battery power:", flightData.BatteryPercentage)
	})
}
