package input

import (
	"github.com/mrechtien/mixgo/config"
)

type InputEvent struct {
	Channel uint8
	Control uint8
	Value   uint8
}

type InputSource interface {
	Setup(config *config.Config, inputEvents chan *InputEvent)
	Reset()
}

type InputSourceCreator func(name string) *InputSource

var inputRegistry = map[string]InputSourceCreator{}

func AddInputSource(inputType string, creator InputSourceCreator) {
	inputRegistry[inputType] = creator
}

func CreateInputSource(inputType string, name string) *InputSource {
	return inputRegistry[inputType](name)
}

func NewInputEvent(control uint8, value uint8) *InputEvent {
	return NewInputEventWithChannel(0, control, value)
}

func NewInputEventWithChannel(channel uint8, control uint8, value uint8) *InputEvent {
	return &InputEvent{
		Channel: channel,
		Control: control,
		Value:   value,
	}
}
