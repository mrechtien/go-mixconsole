package cq

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/mrechtien/mixgo/display"
	"github.com/mrechtien/mixgo/mixer"
)

const (
	MIXER_NAME = "cq"
)

func init() {
	mixer.AddMixer(MIXER_NAME, func(ip string, port uint, displayEvents chan *display.DisplayEvent) *mixer.Mixer {
		return NewMixer(ip, port, displayEvents)
	})
}

type CqMixer struct {
	output        chan []uint8
	displayEvents chan *display.DisplayEvent
}

func NewMixer(ip string, port uint, displayEvents chan *display.DisplayEvent) *mixer.Mixer {
	cqMixer := CqMixer{
		output:        make(chan []uint8),
		displayEvents: displayEvents,
	}
	go sendToMixer(ip, port, cqMixer.output)
	var mixer mixer.Mixer = &cqMixer
	return &mixer
}

func sendToMixer(ip string, port uint, output chan []uint8) {
	for message := range output {
		dialer := net.Dialer{Timeout: (time.Second * 5)}
		serverIpPort := fmt.Sprintf("%s:%d", ip, port)
		connection, err := dialer.Dial("tcp", serverIpPort)
		if err != nil {
			slog.Error("could not connect tcp server ", slog.Any("address", serverIpPort), slog.Any("error", err))
		} else {
			slog.Info("sending message to mixer", slog.Any("int", fmt.Sprintf("%v", message)), slog.String("hex", fmt.Sprintf("% 02X", message)))
			connection.Write(message)
		}
		connection.Close()
	}
}

func (mix *CqMixer) NewMuteGroup(muteChannel uint8) *mixer.MuteGroup {
	var muteGroup mixer.MuteGroup = NewMuteGroup(0x00, muteChannel, mix)
	return &muteGroup
}

func (mix *CqMixer) NewMuteChannel(muteChannel uint8) *mixer.MuteChannel {
	var muteChan mixer.MuteChannel = NewMuteChannel(0x00, muteChannel, mix)
	return &muteChan
}

func (mix *CqMixer) NewTapDelay(fxChannel uint8) *mixer.TapDelay {
	var tapDelay mixer.TapDelay = NewTapDelay(0x00, fxChannel, mix)
	return &tapDelay
}
