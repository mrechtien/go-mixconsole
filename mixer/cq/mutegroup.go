package cq

const (
	MUTE_ON     = 0x01
	MUTE_OFF    = 0x00
	MUTE_GROUPS = 0x00
)

type CqMuteGroup struct {
	midiChannel uint8
	muteChannel uint8
	output      chan []uint8
}

func NewMuteGroup(midiChannel uint8, muteChannel uint8, output chan []uint8) *CqMuteGroup {
	muteGroup := CqMuteGroup{
		midiChannel: midiChannel,
		muteChannel: MUTE_GROUPS + muteChannel,
		output:      output,
	}
	return &muteGroup
}

func (muteGroup *CqMuteGroup) Toggle(onOff bool) {
	message := toMuteGroupMessage(muteGroup.muteChannel, onOff)
	muteGroup.output <- message
}

func toMuteGroupMessage(muteChannel uint8, onOff bool) []uint8 {
	msg := []uint8{0xB0, 0x63, 0x04, 0xB0, 0x62, 0x00, 0xB0, 0x06, 0x00, 0xB0, 0x26, 0x00}
	msg[5] = muteChannel
	if onOff {
		msg[11] = MUTE_ON
	} else {
		msg[11] = MUTE_OFF
	}
	return msg
}
