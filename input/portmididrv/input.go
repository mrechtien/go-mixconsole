package portmididrv

import (
	"log"

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
	input.AddInput(INPUT_TYPE, func(name string) *input.Input {
		return NewInput(name)
	})
}

func NewInput(name string) *input.Input {
	portMidiInput := PortMidiInput{
		inputs: discoverInputs(),
	}
	var input input.Input = &portMidiInput
	return &input
}

func discoverInputs() map[string]portmidi.DeviceID {
	var inputs = make(map[string]portmidi.DeviceID)
	deviceCount := portmidi.CountDevices()
	for index := 0; index < deviceCount; index++ {
		deviceInfo := portmidi.Info(portmidi.DeviceID(index))
		if deviceInfo.IsInputAvailable {
			log.Printf("Input Device '%s'\n", deviceInfo.Name)
			inputs[deviceInfo.Name] = portmidi.DeviceID(index)
		}
		if deviceInfo.IsInputAvailable {
			log.Printf("Output Device '%s'\n", deviceInfo.Name)
		}
	}
	return inputs
}

func (input *PortMidiInput) Setup(config *config.Config, callback func(uint8, uint8, uint8)) func() {
	portmidi.Initialize()
	defer portmidi.Terminate()

	inputName, isPresent := input.inputs[config.Input.Name]
	if !isPresent {
		log.Fatalln("MIDI device not found: %s", inputName)
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
		log.Printf("MIDI Event CC Number [%v] CC Value [%v] Status [%v]", evt.Data1, evt.Data2, evt.Status)
		callback(uint8(evt.Data1), uint8(evt.Data2), uint8(evt.Status))
	}

	// TODO
	return func() {}
}
