package xr

import (
	"log"

	"github.com/hypebeast/go-osc/osc"
	"github.com/mrechtien/mixgo/mixer"
)

const (
	MIXER_NAME = "xr"
)

type XRMixer struct {
	output chan osc.Message
}

func init() {
	mixer.AddMixer(MIXER_NAME, func(ip string, port uint) *mixer.Mixer {
		return NewMixer(ip, port)
	})
}

func NewMixer(ip string, port uint) *mixer.Mixer {
	xrMixer := XRMixer{
		output: make(chan osc.Message),
	}
	go sendToMixer(ip, port, xrMixer.output)
	var mixer mixer.Mixer = &xrMixer
	return &mixer
}

func sendToMixer(ip string, port uint, output chan osc.Message) {
	client := osc.NewClient(ip, int(port))
	for message := range output {
		log.Printf("Sending message to mixer: %v\n", message)
		if err := client.Send(&message); err != nil {
			log.Println("Error while sending message: ", err)
		}
	}
}

func (mix *XRMixer) NewMuteGroup(muteChannel uint8) *mixer.MuteGroup {
	var muteGroup mixer.MuteGroup = NewMuteGroup(muteChannel, mix.output)
	return &muteGroup
}

func (mix *XRMixer) NewMuteChannel(muteChannel uint8) *mixer.MuteChannel {
	var muteGroup mixer.MuteChannel = NewMuteChannel(muteChannel, mix.output)
	return &muteGroup
}

func (mix *XRMixer) NewTapDelay(fxChannel uint8) *mixer.TapDelay {
	var tapDelay mixer.TapDelay = NewTapDelay(fxChannel, mix.output)
	return &tapDelay
}
