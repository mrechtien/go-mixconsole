package cq

import "github.com/mrechtien/mixgo/display"

const (
	MUTE_ON     = 0x01
	MUTE_OFF    = 0x00
	MUTE_GROUPS = 0x00
)

type CqMuteGroup struct {
	midiChannel uint8
	muteChannel uint8
	mixer       *CqMixer
}

func NewMuteGroup(midiChannel uint8, muteChannel uint8, mixer *CqMixer) *CqMuteGroup {
	muteGroup := CqMuteGroup{
		midiChannel: midiChannel,
		muteChannel: MUTE_GROUPS + muteChannel,
		mixer:       mixer,
	}
	return &muteGroup
}

func (control *CqMuteGroup) Toggle(onOff bool) {
	message := toMuteGroupMessage(control.muteChannel, onOff)
	control.mixer.output <- message
	if onOff {
		control.mixer.displayEvents <- display.CreateDisplayEvent(control.muteChannel, uint8(1))
	} else {
		control.mixer.displayEvents <- display.CreateDisplayEvent(control.muteChannel, uint8(0))
	}
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
