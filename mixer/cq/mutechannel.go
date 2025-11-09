package cq

import "github.com/mrechtien/mixgo/display"

type CqMuteChannel struct {
	midiChannel uint8
	muteChannel uint8
	mixer       *CqMixer
}

func NewMuteChannel(midiChannel uint8, muteChannel uint8, mixer *CqMixer) *CqMuteChannel {
	return &CqMuteChannel{
		midiChannel: midiChannel,
		muteChannel: muteChannel,
		mixer:       mixer,
	}
}

func (control *CqMuteChannel) Toggle(onOff bool) {
	message := toMuteChannelMessage(control.muteChannel, onOff)
	control.mixer.output <- message
	if onOff {
		control.mixer.displayEvents <- display.CreateDisplayEvent(control.muteChannel, uint8(1))
	} else {
		control.mixer.displayEvents <- display.CreateDisplayEvent(control.muteChannel, uint8(0))
	}
}

func toMuteChannelMessage(muteChannel uint8, onOff bool) []uint8 {
	msg := []uint8{0xB0, 0x63, 0x00, 0xB0, 0x62, 0x00, 0xB0, 0x06, 0x00, 0xB0, 0x26, 0x00}
	msg[5] = muteChannel
	if onOff {
		msg[11] = MUTE_ON
	} else {
		msg[11] = MUTE_OFF
	}
	return msg
}
