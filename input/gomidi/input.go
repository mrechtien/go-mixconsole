package gomidi

import (
	"log"

	"github.com/mrechtien/mixgo/config"

	"gitlab.com/gomidi/midi/v2"
	//_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

func printMidiDevices() {
	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	log.Printf("MIDI IN Ports\n")
	for i, port := range midi.GetInPorts() {
		log.Printf("no: %v %q\n", i, port)
	}
	log.Printf("\n\nMIDI OUT Ports\n")
	for i, port := range midi.GetOutPorts() {
		log.Printf("no: %v %q\n", i, port)
	}
	log.Printf("\n\n")
}

func Setup(config *config.Config, callback func(uint8, uint8, uint8)) func() {

	// TODO
	/*
		stop := input.SetupAndHandleMidi(&cfg, func(ch, status, val uint8) {
			key := midiToKey(ch, status)
			callback := callbacks[key]
			if callback == nil {
				log.Printf("Unmapped MIDI control change: %+v\n", midi.ControlChange(ch, status, val))
				return
			}
			callback.(func(ch, status, val uint8))(ch, status, val)
		})

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		log.Println("MixGo is up and running!")

		signal := <-signalChan
		log.Printf("Exitting on signal: %d\n", signal)
		stop()
		midi.CloseDriver()
		log.Println("Done.")
	*/

	printMidiDevices()

	in, err := midi.FindInPort(config.Input.Name)
	if err != nil {
		log.Fatalln("can't find given MIDI input device")
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []uint8
		var ch, status, val uint8
		switch {
		case msg.GetControlChange(&ch, &status, &val):
			log.Printf("Received MIDI control change: channel %02X, status %02X, value %02X\n", ch, status, val)
			callback(ch, status, val)
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
			log.Printf("Unmapped MIDI event:  %+v\n", bt)
		}
	}, midi.UseSysEx())

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}

	return stop
}
