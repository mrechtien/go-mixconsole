package portmididrv

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/input"
	"github.com/rakyll/portmidi"
)

const (
	INPUT_TYPE = "portmidi"
)

type PortMidiInput struct {
	inputs map[string]portmidi.DeviceID
}

func init() {
	input.AddInputSource(INPUT_TYPE, func(name string) *input.InputSource {
		return NewInputSource(name)
	})
}

func NewInputSource(name string) *input.InputSource {
	portMidiInput := PortMidiInput{
		inputs: discoverInputs(),
	}
	var input input.InputSource = &portMidiInput
	return &input
}

func discoverInputs() map[string]portmidi.DeviceID {
	var inputs = make(map[string]portmidi.DeviceID)
	deviceCount := portmidi.CountDevices()
	for index := 0; index < deviceCount; index++ {
		deviceInfo := portmidi.Info(portmidi.DeviceID(index))
		if deviceInfo.IsInputAvailable {
			slog.Info("input device", slog.Any("name", deviceInfo.Name))
			inputs[deviceInfo.Name] = portmidi.DeviceID(index)
		}
		if deviceInfo.IsInputAvailable {
			slog.Info("output device", slog.Any("name", deviceInfo.Name))
		}
	}
	return inputs
}

func toHex(input any) string {
	return fmt.Sprintf("%02X", input)
}

func (portMidiInput *PortMidiInput) Setup(config *config.Config, inputEvents chan *input.InputEvent) {
	portmidi.Initialize()
	defer portmidi.Terminate()

	inputName, isPresent := portMidiInput.inputs[config.Input.Name]
	if !isPresent {
		log.Fatalf("MIDI device not found: %s", inputName)
	}

	in, err := portmidi.NewInputStream(inputName, 1024)
	if err != nil {
		log.Fatal(err)
		log.Fatalln("can't find given MIDI input device")
	}
	defer in.Close()

	// or alternatively listen events
	for evt := range in.Listen() {
		// Data1 = CC Number
		// Data2 = CC Value
		slog.Debug("Received MIDI event", slog.Any("channel", toHex(evt.Data1)), slog.Any("control change", toHex(evt.Data2)), slog.Any("status", toHex(evt.Status)))
		inputEvents <- input.NewInputEvent(uint8(evt.Data1), uint8(evt.Data2))
		//callback(uint8(evt.Data1), uint8(evt.Data2), uint8(evt.Status))
	}
}
