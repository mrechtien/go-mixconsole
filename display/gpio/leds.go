package gpio

import (
	"log"
	"strconv"
	"time"

	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	DISPLAY_TYPE = "gpio"
)

type LedStateDisplay struct {
}

type LedState struct {
	Mapping config.Mapping
	Gpio    gpio.PinIO
}

type LedTemporary struct {
	LedState
}
type LedBiState struct {
	LedState
	state bool
}

func init() {
	display.AddDisplay(DISPLAY_TYPE, func() display.Display {
		return &LedStateDisplay{}
	})
}

func (ledStateDisplay *LedStateDisplay) CreateBiState(mapping config.Mapping) display.BiState {

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	gpio := gpioreg.ByName(strconv.Itoa(int(mapping.Display)))
	led := LedBiState{
		LedState: LedState{
			Mapping: mapping,
			Gpio:    gpio,
		},
		state: false,
	}

	var biState display.BiState = &led
	return biState
}

func (led *LedState) On() {
	led.Gpio.Out(gpio.Level(true))
}

func (led *LedState) Off() {
	led.Gpio.Out(gpio.Level(false))
}

func (ledStateDisplay *LedStateDisplay) CreateTemporary(mapping config.Mapping) display.Temporary {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	gpio := gpioreg.ByName(strconv.Itoa(int(mapping.Display)))
	led := LedTemporary{
		LedState: LedState{
			Mapping: mapping,
			Gpio:    gpio,
		},
	}

	var temporary display.Temporary = &led
	return temporary
}

func (led *LedTemporary) Trigger() {
	go func() {
		led.Gpio.Out(gpio.Level(true))
		time.Sleep(100 * time.Millisecond)
		led.Gpio.Out(gpio.Level(false))
	}()
}
