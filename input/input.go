package input

import "github.com/mrechtien/mixgo/config"

type Input interface {
	Setup(config *config.Config, callback func(uint8, uint8, uint8)) func()
}

var inputRegistry = map[string]interface{}{}

func AddInput(inputType string, creator func(name string) *Input) {
	inputRegistry[inputType] = creator
}

func CreateInput(inputType string, name string) *Input {
	return inputRegistry[inputType].(func(name string) *Input)(name)
}
