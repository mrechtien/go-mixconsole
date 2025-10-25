package qu

const (
	MUTE_ON     = 0x40
	MUTE_OFF    = 0x10
	MUTE_GROUPS = 0x50 // start channel
)

type QuMuteGroup struct {
	midiChannel uint8
	muteChannel uint8
	output      chan []uint8
}

func NewMuteGroup(midiChan uint8, muteChan uint8, output chan []uint8) *QuMuteGroup {
	muteGroup := QuMuteGroup{
		midiChannel: midiChan,
		muteChannel: MUTE_GROUPS + muteChan,
		output:      output,
	}
	return &muteGroup
}

func (muteGroup *QuMuteGroup) Toggle(onOff bool) {
	muteGroup.output <- createMuteGroupMessage(muteGroup.muteChannel, onOff)
}

func createMuteGroupMessage(muteChan uint8, onOff bool) []uint8 {
	msg := []uint8{0x90, 0x00, 0x7F, 0x90, 0x00, 0x40}
	msg[1] = muteChan
	msg[4] = muteChan
	if onOff {
		msg[5] = MUTE_ON
	} else {
		msg[5] = MUTE_OFF
	}
	return msg
}
