package setup

import (
	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/input"
)

func SetupMixGo(cfg *config.Config) {

}

func Bootstrap() *config.Config {

	cfg := config.LoadBootstrapConfig()

	inputSource := *input.CreateInputSource(cfg.Input.Type, cfg.Input.Name)
	inputEvents := make(chan *input.InputEvent)
	inputSource.Setup(cfg, inputEvents)

}
