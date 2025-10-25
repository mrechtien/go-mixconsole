package qu

import "log"

type QuMuteChannel struct {
	midiChannel uint8
	muteChannel uint8
	output      chan []uint8
}

func NewMuteChannel(midiChan uint8, muteChan uint8, output chan []uint8) *QuMuteChannel {
	muteChannel := QuMuteChannel{
		midiChannel: midiChan,
		muteChannel: MUTE_GROUPS + muteChan,
		output:      output,
	}
	return &muteChannel
}

func (muteChannel *QuMuteChannel) Toggle(onOff bool) {
	log.Fatalln("Unsupported method")
}
