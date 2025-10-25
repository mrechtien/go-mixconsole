package xr

import (
	"log"

	"github.com/hypebeast/go-osc/osc"
	"github.com/mrechtien/mixgo/base"
)

type XRMuteChannel struct {
	base.MuteGroup
	muteChannel uint8
	output      chan osc.Message
}

func NewMuteChannel(muteChannel uint8, output chan osc.Message) *XRMuteChannel {
	xrMuteChannel := XRMuteChannel{
		muteChannel: muteChannel,
		output:      output,
	}
	return &xrMuteChannel
}

func (muteGroup *XRMuteChannel) Toggle(onOff bool) {
	log.Fatalln("Unsupported method")
}
