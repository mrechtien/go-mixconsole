package gpio

import (
	"log"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"

	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/input"
)

const (
	INPUT_TYPE = "gpio"
)

type GPIOInput struct {
	Control uint8
}

func init() {
	input.AddInputSource(INPUT_TYPE, func(name string) *input.InputSource {
		return NewInputSource(name)
	})
}

func NewInputSource(name string) *input.InputSource {
	gpioInput := GPIOInput{}
	var input input.InputSource = &gpioInput
	return &input
}

func initGPIO() {

	slog.Info("initialising periph.io system")

	// load and init drivers
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	allPins := gpioreg.All()
	regex := regexp.MustCompile("^GPIO[0-9]([0-9])?$")
	slog.Info("available [0-99] GPIO pins:")
	for _, pin := range allPins {
		isMatch := regex.MatchString(pin.Name())
		if isMatch {
			slog.Info("found GPIO pin", slog.Any("pin", pin))
		}
	}
}

func listenToGpio(mapping config.Mapping, inputEvents chan *input.InputEvent) {

	// Lookup a pin by its number:
	control := strconv.Itoa(int(mapping.Control))
	port := gpioreg.ByName(control)
	if port == nil {
		log.Fatalf("failed to find GPIO [%v]", port)
	}

	slog.Debug("listen to gpio", slog.Any("port", port), slog.Any("function", port.Function()))

	// Set it as input, with an internal pull down resistor:
	if err := port.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}

	// Wait for edges as detected by the hardware, and print the value read:
	var lastState gpio.Level
	var toggle bool
	for {
		port.WaitForEdge(-1)
		state := port.Read()
		if lastState != state {
			if state {
				toggle = !toggle
				slog.Debug("handling gpio event", slog.Any("state", state), slog.Any("lastState", lastState), slog.Any("toggle", toggle))
				if toggle {
					inputEvents <- input.NewInputEvent(mapping.Control, 1)
				} else {
					inputEvents <- input.NewInputEvent(mapping.Control, 0)
				}
			}
			lastState = state
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (input *GPIOInput) Setup(config *config.Config, inputEvents chan *input.InputEvent) {

	// init
	initGPIO()

	// create GPIO inputs
	for _, mapping := range config.Mappings {
		go listenToGpio(mapping, inputEvents)
	}
}
