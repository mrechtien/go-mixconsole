package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func switchLed(readPin chan bool) {

	p := gpioreg.ByName("26")

	for state := range readPin {
		log.Printf("-> %v\n", state)
		p.Out(gpio.Level(state))
	}
}

func listenToGpio(readPin chan bool) {

	// Lookup a pin by its number:
	p := gpioreg.ByName("17")
	if p == nil {
		log.Fatal("Failed to find GPIO26")
	}

	log.Printf("%s: %s\n", p, p.Function())

	// Set it as input, with an internal pull down resistor:
	if err := p.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}

	// Wait for edges as detected by the hardware, and print the value read:
	var lastState gpio.Level
	for {
		p.WaitForEdge(-1)
		state := p.Read()
		if lastState != state {
			log.Printf("-> %v\n", state)
			readPin <- bool(state)
			lastState = state
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	allPins := gpioreg.All()
	regex := regexp.MustCompile("^GPIO[0-9]([0-9])?$")
	log.Println("Available pins:")
	for _, pin := range allPins {
		isMatch := regex.MatchString(pin.Name())
		if isMatch {
			log.Printf("%v", pin)
		}
	}

	readPin := make(chan bool)
	go listenToGpio(readPin)
	go switchLed(readPin)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Up and running ....!")

	signal := <-signalChan
	log.Printf("Exitting on signal: %d\n", signal)

}
