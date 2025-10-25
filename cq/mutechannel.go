package cq

type CqMuteChannel struct {
	midiChannel uint8
	muteChannel uint8
	output      chan []uint8
}

func NewMuteChannel(midiChannel uint8, muteChannel uint8, output chan []uint8) *CqMuteChannel {
	cqMuteChannel := CqMuteChannel{
		midiChannel: midiChannel,
		muteChannel: muteChannel,
		output:      output,
	}
	return &cqMuteChannel
}

func (muteChannel *CqMuteChannel) Toggle(onOff bool) {
	message := toMuteChannelMessage(muteChannel.muteChannel, onOff)
	muteChannel.output <- message
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
