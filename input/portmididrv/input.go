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

var defaultChannel uint8

func (portMidiInput *PortMidiInput) Setup(config *config.Config, inputEvents chan *input.InputEvent) {
	portmidi.Initialize()

	defaultChannel = config.Input.Channel
	if defaultChannel < 1 {
		defaultChannel = 1
	}
	inputName, isPresent := portMidiInput.inputs[config.Input.Name]
	if !isPresent {
		log.Fatalf("MIDI device not found: %v", inputName)
	}

	go listenToMidiIn(inputName, inputEvents)
}

func (portMidiInput *PortMidiInput) Reset() {
	// TODO XXX
}

func listenToMidiIn(inputName portmidi.DeviceID, inputEvents chan *input.InputEvent) {
	midiIn, err := portmidi.NewInputStream(inputName, 1024)
	if err != nil {
		log.Fatal(err)
		log.Fatalln("can't find given MIDI input device")
	}
	defer midiIn.Close()

	// or alternatively listen events
	for evt := range midiIn.Listen() {
		// Data1 = CC Number
		// Data2 = CC Value
		slog.Debug("received MIDI event", slog.Any("channel", evt.Data1), slog.Any("control change", evt.Data2), slog.Any("status", evt.Status))
		inputEvents <- input.NewInputEventWithChannel(getChannel(uint8(evt.Status)), uint8(evt.Data1), uint8(evt.Data2))
	}
	portmidi.Terminate()
}

func getChannel(status uint8) uint8 {
	ch := defaultChannel
	cmd := status & 0xF0 // mask off all but top 4 bits
	if cmd >= 0x80 && cmd <= 0xE0 {
		// find the channel by masking off all but the low 4 bits
		ch = (status & 0x0F) + 1
	}
	return ch
}
