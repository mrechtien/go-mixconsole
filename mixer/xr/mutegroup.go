package xr

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
	"github.com/mrechtien/mixgo/mixer"
)

const (
	MUTE_ON  = int32(1)
	MUTE_OFF = int32(0)
)

type XRMuteGroup struct {
	mixer.MuteGroup
	muteChannel uint8
	output      chan osc.Message
}

func NewMuteGroup(muteChannel uint8, output chan osc.Message) *XRMuteGroup {
	muteGroup := XRMuteGroup{
		muteChannel: muteChannel,
		output:      output,
	}
	return &muteGroup
}

func (muteGroup *XRMuteGroup) Toggle(onOff bool) {
	value := MUTE_OFF
	if onOff {
		value = MUTE_ON
	}

	message := osc.NewMessage(fmt.Sprintf("/config/mute/%d", muteGroup.muteChannel), value)
	muteGroup.output <- *message
}
