package cq

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrechtien/mixgo/base"
)

const (
	MIXER_NAME = "cq"
)

func init() {
	base.AddMixer(MIXER_NAME, func(ip string, port uint) *base.Mixer {
		return NewMixer(ip, port)
	})
}

type CqMixer struct {
	output chan []uint8
}

func NewMixer(ip string, port uint) *base.Mixer {
	cqMixer := CqMixer{
		output: make(chan []uint8),
	}
	go sendToMixer(ip, port, cqMixer.output)
	var mixer base.Mixer = &cqMixer
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

func (mixer *CqMixer) NewMuteGroup(muteChan uint8) *base.MuteGroup {
	var muteGroup base.MuteGroup = NewMuteGroup(0x00, muteChan, mixer.output)
	return &muteGroup
}

func (mixer *CqMixer) NewMuteChannel(muteChan uint8) *base.MuteChannel {
	var muteChannel base.MuteChannel = NewMuteChannel(0x00, muteChan, mixer.output)
	return &muteChannel
}

func (mixer *CqMixer) NewTapDelay(fxChan uint8) *base.TapDelay {
	var tapDelay base.TapDelay = NewTapDelay(0x00, fxChan, mixer.output)
	return &tapDelay
}
