package mixer

import "github.com/mrechtien/mixgo/display"

type Mixer interface {
	NewMuteGroup(muteChan uint8) *MuteGroup
	NewMuteChannel(muteChan uint8) *MuteChannel
	NewTapDelay(fxChan uint8) *TapDelay
}
type MixerCreator func(ip string, port uint, displayEvents chan *display.DisplayEvent) *Mixer

var mixerRegistry = map[string]MixerCreator{}

func AddMixer(name string, creator MixerCreator) {
	mixerRegistry[name] = creator
}

func CreateMixer(name string, ip string, port uint, displayEvents chan *display.DisplayEvent) *Mixer {
	return mixerRegistry[name](ip, port, displayEvents)
}
