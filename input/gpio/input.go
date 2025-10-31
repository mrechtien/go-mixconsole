package gpio

import (
	"log"

	"github.com/mrechtien/mixgo/config"
)

func printConfiguredPorts() {
	log.Printf("Configured ports\n")

	log.Printf("\n\n")
}

func Setup(config *config.Config, callback func(uint8, uint8, uint8)) func() {

	printConfiguredPorts()

	/*

		in, err := midi.FindInPort(config.Input.Name)
		if err != nil {
			log.Fatalln("can't find given GPIO input")
		}
	*/
	// add listeners

	var err interface{}
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}

	return func() {}
}
