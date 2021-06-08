package main

import (
	"github.com/sujayramesh/SMS/api"

	"flag"
	"fmt"
	"sync"
	"time"
)

var Wg sync.WaitGroup

func handleEvent(periodicTicker *time.Ticker) {
	for {
		select {
		case t := <-periodicTicker.C:
			{
				fmt.Println("Periodic timer expired: ", t)
				api.TriggerStaleSessionCleanup()
			}
		}
	}
}

func InitServerDependencies() {

	fmt.Println("Enter StartServer")

	e := api.NewEcho()
	e.Start(":7443")
}

func main() {

	fmt.Println("Starting service!!")

	periodicTimerDuration := flag.Int("PeriodicTimer", 300, "Configurable periodic timer value")
	if *periodicTimerDuration < 0 || *periodicTimerDuration > 600 {
		fmt.Println("Periodic timer value passed: ", *periodicTimerDuration, ", defaulting to 300s.")
		*periodicTimerDuration = 300
	}
	fmt.Println("Periodic timer duration: ", *periodicTimerDuration)

	Wg.Add(1)

	// Init server credentials
	go InitServerDependencies()

	// Clean up session stored in memory.
	defer api.TriggerExitActions()

	fmt.Println("Service listening on port")

	// Init periodic timer
	periodicTicker := time.NewTicker(time.Duration(*periodicTimerDuration) * time.Second)
	defer periodicTicker.Stop()

	handleEvent(periodicTicker)

	Wg.Wait()

}
