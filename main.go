package main

import (
	_ "github.com/mrechtien/mixgo/display/gpio"
	_ "github.com/mrechtien/mixgo/input/gomidi"
	_ "github.com/mrechtien/mixgo/input/gpio"
	_ "github.com/mrechtien/mixgo/input/portmididrv"
	_ "github.com/mrechtien/mixgo/mixer/cq"
	"github.com/mrechtien/mixgo/setup"
)

func main() {

	setup.SetupLogging()

	cfg := setup.Bootstrap()

	setup.SetupAndRun(cfg)
}
