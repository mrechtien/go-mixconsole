package base

type Mixer interface {
	NewMuteGroup(muteChan uint8) *MuteGroup
	NewMuteChannel(muteChan uint8) *MuteChannel
	NewTapDelay(fxChan uint8) *TapDelay
}

var mixerRegistry = map[string]interface{}{}

func AddMixer(name string, creator func(ip string, port uint) *Mixer) {
	mixerRegistry[name] = creator
}

func CreateMixer(name string, ip string, port uint) *Mixer {
	return mixerRegistry[name].(func(ip string, port uint) *Mixer)(ip, port)
}
