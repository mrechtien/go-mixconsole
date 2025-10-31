package qu

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrechtien/mixgo/mixer"
)

const (
	MIXER_NAME = "qu"
)

func init() {
	mixer.AddMixer(MIXER_NAME, func(ip string, port uint) *mixer.Mixer {
		return NewMixer(ip, port)
	})
}

type QuMixer struct {
	output chan []uint8
}

func NewMixer(ip string, port uint) *mixer.Mixer {
	quMixer := QuMixer{
		output: make(chan []uint8),
	}
	go sendToMixer(ip, port, quMixer.output)
	var mixer mixer.Mixer = &quMixer
	return &mixer
}

func sendToMixer(ip string, port uint, output chan []uint8) {
	for message := range output {
		dialer := net.Dialer{Timeout: (time.Second * 5)}
		connection, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			log.Printf("Could not connect to TCP server: %s", err)
		} else {
			log.Printf("Sending message to mixer: %v\n", message)
			connection.Write(message)
		}
		connection.Close()
	}
}

func (mix *QuMixer) NewMuteGroup(muteChannel uint8) *mixer.MuteGroup {
	var muteGroup mixer.MuteGroup = NewMuteGroup(0x00, muteChannel, mix.output)
	return &muteGroup
}

func (mix *QuMixer) NewMuteChannel(muteChannel uint8) *mixer.MuteChannel {
	var muteChan mixer.MuteChannel = NewMuteGroup(0x00, muteChannel, mix.output)
	return &muteChan
}

func (mix *QuMixer) NewTapDelay(fxChannel uint8) *mixer.TapDelay {
	var tapDelay mixer.TapDelay = NewTapDelay(0x00, fxChannel, mix.output)
	return &tapDelay
}
