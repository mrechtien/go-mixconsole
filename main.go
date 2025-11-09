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

type InputEventConsumer func(cc uint8, val uint8, status uint8)

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

	// configuration
	var cfg config.Config
	if len(os.Args) == 2 {
		configPath := os.Args[1]
		cfg = config.ReadConfig(configPath)
	}

	// wifi
	if len(cfg.Wifi.SSID) > 0 {
		//setup.ConnectWifi(&cfg.Wifi)
	}

	// display
	displayEvents := make(chan *display.DisplayEvent)
	display.SetupDisplay(&cfg.Display, cfg.Mappings, displayEvents)

	// mixer
	mix := *mixer.CreateMixer(cfg.Output.Name, cfg.Output.Ip, cfg.Output.Port, displayEvents)

	// create callbacks for trigger mapping
	callbacks := map[string]InputEventConsumer{}
	for _, mapping := range cfg.Mappings {
		key := midiToKey(cfg.Input.Channel, mapping.Control)
		slog.Debug("create callback", slog.Any("key", key))
		switch mapping.Name {
		case mixer.MUTE_GROUP:
			muteGroup := *mix.NewMuteGroup(mapping.Target)
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				slog.Debug("callback MuteGroup", slog.Any("value", val))
				muteGroup.Toggle(val == mapping.ValueOn)
			}
		case mixer.MUTE_CHANNEL:
			muteChannel := *mix.NewMuteChannel(mapping.Target)
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				slog.Debug("callback MuteChannel", slog.Any("value", val))
				muteChannel.Toggle(val == mapping.ValueOn)
			}
		case mixer.TAP_DELAY:
			tapDelay := *mix.NewTapDelay(mapping.Target)
			slog.Debug("callback TapDelay")
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				tapDelay.Trigger()
			}
		default:
			log.Fatalln("Invalid mapping name in config: ", mapping.Name)
		}
	}

	// setup input
	inputSource := *input.CreateInputSource(cfg.Input.Type, cfg.Input.Name)
	inputEvents := make(chan *input.InputEvent)

	inputSource.Setup(&cfg, inputEvents)

	for event := range inputEvents {

		slog.Debug("received input event", slog.Any("value", event))

		// TODO fixed input channel should come from input itself
		ch := cfg.Input.Channel
		// TODO midi channel is implementation specific should be internally (before)
		var status uint8     // TODO remove and replace with event.channel use
		cmd := status & 0xF0 // mask off all but top 4 bits
		if cmd >= 0x80 && cmd <= 0xE0 {
			// it's a voice message
			// find the channel by masking off all but the low 4 bits
			ch = (status & 0x0F) + 1
		}

		key := midiToKey(ch, event.Control)
		callback := callbacks[key]
		if callback == nil {
			slog.Warn("unmapped MIDI event", slog.Any("control", event.Control), slog.Any("value", event.Value), slog.Any("status", status))
			return
		}
		slog.Info("mapped MIDI event", slog.Any("control", event.Control), slog.Any("value", event.Value), slog.Any("status", status))
		callback(event.Channel, event.Value, status)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	slog.Info("MixGo is up and running!")
	signal := <-signalChan

	slog.Info("exitting", slog.Any("signal", signal))
	slog.Info("done.")
}
