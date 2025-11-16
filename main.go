package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrechtien/mixgo/config"
	"github.com/mrechtien/mixgo/display"
	_ "github.com/mrechtien/mixgo/display/gpio"
	"github.com/mrechtien/mixgo/input"
	_ "github.com/mrechtien/mixgo/input/gomidi"
	_ "github.com/mrechtien/mixgo/input/gpio"
	_ "github.com/mrechtien/mixgo/input/portmididrv"
	"github.com/mrechtien/mixgo/mixer"
	_ "github.com/mrechtien/mixgo/mixer/cq"
)

/**
 * creates a key for callback / action lookup using
 * MIDI channel and control change value
 */
func midiToKey(ch uint8, cc uint8) string {
	return fmt.Sprintf("%02X%02X", ch, cc)
}

type InputEventConsumer func(value uint8)

func setupLogging() {
	// logging
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)
}



func main() {
	setupLogging()

	cfg := bootstrap()

	// setup wifi
	if len(cfg.Wifi.SSID) > 0 {
		//setup.ConnectWifi(&cfg.Wifi)
	}

	// setup display
	displayEvents := make(chan *display.DisplayEvent)
	display.SetupDisplay(&cfg.Display, cfg.Mappings, displayEvents)

	// setup mixer
	mix := *mixer.CreateMixer(cfg.Output.Name, cfg.Output.Ip, cfg.Output.Port, displayEvents)

	// setup callbacks for input event mapping
	callbacks := map[string]InputEventConsumer{}
	for _, mapping := range cfg.Mappings {
		key := midiToKey(cfg.Input.Channel, mapping.Control)
		slog.Debug("create callback", slog.Any("key", key))
		switch mapping.Name {
		case mixer.MUTE_GROUP:
			muteGroup := *mix.NewMuteGroup(mapping.Target)
			callbacks[key] = func(value uint8) {
				slog.Debug("callback MuteGroup", slog.Any("value", value))
				muteGroup.Toggle(value == mapping.ValueOn)
			}
		case mixer.MUTE_CHANNEL:
			muteChannel := *mix.NewMuteChannel(mapping.Target)
			callbacks[key] = func(value uint8) {
				slog.Debug("callback MuteChannel", slog.Any("value", value))
				muteChannel.Toggle(value == mapping.ValueOn)
			}
		case mixer.TAP_DELAY:
			tapDelay := *mix.NewTapDelay(mapping.Target)
			slog.Debug("callback TapDelay")
			callbacks[key] = func(value uint8) {
				tapDelay.Trigger()
			}
		default:
			log.Fatalln("Invalid mapping name in config: ", mapping.Name)
		}
	}

	// setup input and handle events
	inputSource := *input.CreateInputSource(cfg.Input.Type, cfg.Input.Name)
	inputEvents := make(chan *input.InputEvent)
	inputSource.Setup(&cfg, inputEvents)
	for event := range inputEvents {
		slog.Debug("received input event", slog.Any("value", event))
		key := midiToKey(event.Channel, event.Control)
		callback := callbacks[key]
		if callback == nil {
			slog.Warn("unmapped MIDI event", slog.Any("control", event.Control), slog.Any("value", event.Value))
			return
		}
		slog.Info("mapped MIDI event", slog.Any("control", event.Control), slog.Any("value", event.Value))
		go callback(event.Value)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	slog.Info("MixGo is up and running!")
	signal := <-signalChan

	slog.Info("exitting", slog.Any("signal", signal))
	slog.Info("done.")
}
