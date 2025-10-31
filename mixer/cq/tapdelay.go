package cq

import "github.com/mrechtien/mixgo/mixer"

type CqTapDelay struct {
	mixer.BaseTapDelay
	midiChannel uint8
	fxChannel   uint8
	output      chan []uint8
}

// channel is the mixer channel (FX) to trigger the tap delay on
func NewTapDelay(midiChannel uint8, fxChannel uint8, output chan []uint8) *CqTapDelay {
	tapDelay := CqTapDelay{
		BaseTapDelay: mixer.BaseTapDelay{
			LastTriggered: 0,
			Tapping:       []int64{},
		},
		midiChannel: midiChannel,
		fxChannel:   fxChannel,
		output:      output,
	}
	return &tapDelay
}

/**
 *
 */
func (tapDelay *CqTapDelay) Trigger() {
	tapDelay.output <- createSoftKeyDownMessage()
}

func createSoftKeyDownMessage() []uint8 {
	// softkey release is not neede it seems
	// 0x80, 0x32, 0x00
	message := []uint8{0x90, 0x32, 0x7F}
	return message
}
