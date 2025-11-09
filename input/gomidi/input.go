package gomidi

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/input"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	//_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

func printMidiDevices() {
	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	slog.Info("listing MIDI inputs")
	for i, port := range midi.GetInPorts() {
		slog.Info("available MIDI input port", slog.Any("idx", i), slog.Any("port", port))
	}
	slog.Info("listing MIDI outputs")
	for i, port := range midi.GetOutPorts() {
		slog.Info("available MIDI output port", slog.Any("idx", i), slog.Any("port", port))
	}
}

func toHex(input any) string {
	return fmt.Sprintf("%02X", input)
}

func Setup(config *config.Config, inputEvents chan *input.InputEvent) {
	printMidiDevices()

	midiIn, err := midi.FindInPort(config.Input.Name)
	if err != nil {
		log.Fatalln("can't find given MIDI input device")
	}

	go listenToMidiIn(midiIn, inputEvents)
}

func listenToMidiIn(midiIn drivers.In, inputEvents chan *input.InputEvent) {
	stop, err := midi.ListenTo(midiIn, func(msg midi.Message, timestampms int32) {
		var bt []uint8
		var ch, cc, val uint8
		switch {
		case msg.GetControlChange(&ch, &cc, &val):
			slog.Debug("received MIDI control change", slog.Any("channel", toHex(ch)), slog.Any("control change", toHex(cc)), slog.Any("value", toHex(val)))
			// TODO status => val OR val => status
			inputEvents <- &input.InputEvent{
				Channel: ch,
				Control: cc,
				Value:   val,
			}

			//callback(ch, status, val)
		/*
			case msg.GetSysEx(&bt):
				log.Printf("got sysex: %X\n", bt)
			case msg.GetNoteStart(&ch, &status, &val):
				log.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(status), ch, val)
			case msg.GetNoteEnd(&ch, &status):
				log.Printf("ending note %s on channel %v\n", midi.Note(status), ch)
		*/
		default:
			msg.GetSysEx(&bt)
			slog.Warn("unmapped MIDI event", slog.Any("bt", bt))
		}
	}, midi.UseSysEx())

	if err != nil {
		slog.Error("GoMidi error", slog.Any("error", err))
	}
	defer stop()
}
