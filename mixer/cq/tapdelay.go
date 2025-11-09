package cq

import (
	"log/slog"
	"time"

	"github.com/mrechtien/mixgo/display"
	"github.com/mrechtien/mixgo/mixer"
)

const (
	MAX_DELAY_MILLISECONDS = 2400
	MIN_DELAY_TAPS         = 1
)

type CqTapDelay struct {
	mixer.BaseTapDelay
	midiChannel uint8
	fxChannel   uint8
	mixer       *CqMixer
	ticker      *time.Ticker
}

// channel is the mixer channel (FX) to trigger the tap delay on
func NewTapDelay(midiChannel uint8, fxChannel uint8, cqMixer *CqMixer) *CqTapDelay {
	tapDelay := CqTapDelay{
		BaseTapDelay: mixer.BaseTapDelay{
			LastTriggered: 0,
			Tapping:       []int64{},
		},
		midiChannel: midiChannel,
		fxChannel:   fxChannel,
		mixer:       cqMixer,
	}
	return &tapDelay
}

/**
 *
 */
func (control *CqTapDelay) Trigger() {
	control.mixer.output <- createSoftKeyDownMessage()
	control.mixer.displayEvents <- display.CreateDisplayEvent(control.fxChannel, uint8(1))
	delay := mixer.CalculateTapTempo(&control.BaseTapDelay, MAX_DELAY_MILLISECONDS, MIN_DELAY_TAPS)

	if delay > 0 {
		go setupAndRunTicker(control, delay)
	}
}

func createSoftKeyDownMessage() []uint8 {
	// softkey release is not needed it seems
	// 0x80, 0x32, 0x00
	message := []uint8{0x90, 0x32, 0x7F}
	return message
}

func setupAndRunTicker(control *CqTapDelay, delay int) {
	if control.ticker != nil {
		control.ticker.Stop()
	}
	slog.Debug("CqTapDelay", slog.Int("delay", delay))
	control.ticker = time.NewTicker(time.Duration(delay) * time.Millisecond)
	defer control.ticker.Stop()
	for {
		<-control.ticker.C
		control.mixer.displayEvents <- display.CreateDisplayEvent(control.fxChannel, uint8(1))
	}
}
