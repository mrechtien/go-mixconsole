package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrechtien/mixgo/base"
	"github.com/mrechtien/mixgo/config"
	_ "github.com/mrechtien/mixgo/cq"
	"github.com/mrechtien/mixgo/input"
	_ "github.com/mrechtien/mixgo/input/portmididrv"
)

/**
 * creates a key for callback / action lookup using
 * MIDI channel and control change value
 */
func midiToKey(ch uint8, cc uint8) string {
	return fmt.Sprintf("%02X%02X", ch, cc)
}

/**
 *
 */
func main() {
	var cfg config.Config
	if len(os.Args) == 2 {
		configPath := os.Args[1]
		cfg = config.ReadConfig(configPath)
	}

	// setup mixer
	mixer := *base.CreateMixer(cfg.Output.Name, cfg.Output.Ip, cfg.Output.Port)

	// create callbacks for trigger mapping
	callbacks := map[string]interface{}{}
	for _, mapping := range cfg.Mappings {
		key := midiToKey(cfg.Input.Channel, mapping.CC)
		switch mapping.Name {
		case base.MUTE_GROUP:
			muteGroup := *mixer.NewMuteGroup(mapping.Target)
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				muteGroup.Toggle(val == mapping.ValueOn)
			}
		case base.MUTE_CHANNEL:
			muteChannel := *mixer.NewMuteChannel(mapping.Target)
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				muteChannel.Toggle(val == mapping.ValueOn)
			}
		case base.TAP_DELAY:
			tapDelay := *mixer.NewTapDelay(mapping.Target)
			callbacks[key] = func(cc uint8, val uint8, status uint8) {
				tapDelay.Trigger()
			}
		default:
			log.Fatalln("Invalid mapping name in config: ", mapping.Name)
		}
	}

	// setup input
	input := *input.CreateInput(cfg.Input.Type, cfg.Input.Name)

	input.Setup(&cfg, func(cc uint8, val uint8, status uint8) {
		ch := cfg.Input.Channel
		cmd := status & 0xF0 // mask off all but top 4 bits
		if cmd >= 0x80 && cmd <= 0xE0 {
			// it's a voice message
			// find the channel by masking off all but the low 4 bits
			ch = (status & 0x0F) + 1
		}

		key := midiToKey(ch, cc)
		callback := callbacks[key]
		if callback == nil {
			log.Printf("Unmapped MIDI Event CC Number [%v] CC Value [%v] Status [%v]", cc, val, status)
			return
		}
		callback.(func(uint8, uint8, uint8))(cc, val, status)
	})

	/*
		callback := callbacks["0102"]
		go func() {
			for i := range 2 {
				callback.(func(ch, status, val uint8))(0, 1, 0x01)
				// Calling Sleep method
				time.Sleep(2000 * time.Millisecond)
				fmt.Println("xxx ", i)
			}
		}()
	*/

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	log.Println("MixGo is up and running!")
	signal := <-signalChan
	log.Printf("Exitting on signal: %d\n", signal)
	log.Println("Done.")
}
