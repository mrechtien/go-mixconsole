package cq

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrechtien/mixgo/mixer"
)

const (
	MIXER_NAME = "cq"
)

func init() {
	mixer.AddMixer(MIXER_NAME, func(ip string, port uint) *mixer.Mixer {
		return NewMixer(ip, port)
	})
}

type CqMixer struct {
	output chan []uint8
}

func NewMixer(ip string, port uint) *mixer.Mixer {
	cqMixer := CqMixer{
		output: make(chan []uint8),
	}
	go sendToMixer(ip, port, cqMixer.output)
	var mixer mixer.Mixer = &cqMixer
	return &mixer
}

func sendToMixer(ip string, port uint, output chan []uint8) {
	for message := range output {
		dialer := net.Dialer{Timeout: (time.Second * 5)}
		connection, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			log.Printf("Could not connect to TCP server: %s", err)
		} else {
			log.Printf("Sending message to mixer: %v => % 02X\n", message, message)
			connection.Write(message)
		}
		connection.Close()
	}
}

func (mix *CqMixer) NewMuteGroup(muteChannel uint8) *mixer.MuteGroup {
	var muteGroup mixer.MuteGroup = NewMuteGroup(0x00, muteChannel, mix.output)
	return &muteGroup
}

func (mix *CqMixer) NewMuteChannel(muteChannel uint8) *mixer.MuteChannel {
	var muteChan mixer.MuteChannel = NewMuteChannel(0x00, muteChannel, mix.output)
	return &muteChan
}

func (mix *CqMixer) NewTapDelay(fxChannel uint8) *mixer.TapDelay {
	var tapDelay mixer.TapDelay = NewTapDelay(0x00, fxChannel, mix.output)
	return &tapDelay
}
