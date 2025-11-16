package config

import (
	"log"
	"os"

	"github.com/ghodss/yaml"
)

const (
	MUTE_GROUP   = "MuteGroup"
	MUTE_CHANNEL = "MuteChannel"
	TAP_DELAY    = "TapDelay"
)

type Input struct {
	Type    string
	Name    string
	Channel uint8
}

type Output struct {
	Name string
	Ip   string
	Port uint
}

type Mapping struct {
	Name    string
	Control uint8
	Target  uint8
	ValueOn uint8
	Display uint8
	Config  string
}

type Wifi struct {
	SSID   string
	Passwd string
	Hidden bool
}

type Display struct {
	Name string
}

type Config struct {
	Input    Input
	Output   Output
	Mappings []Mapping
	Wifi     Wifi
	Display  Display
}

func NewMapping() *Mapping {
	return &Mapping{
		ValueOn: 1,
	}
}

func ReadConfigFromArgs() *Config {
	// init and read configuration
	var cfg Config
	if len(os.Args) == 2 {
		configPath := os.Args[1]
		cfg = ReadConfig(configPath)
	}
	return &cfg
}

func ReadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("config file not readable, path: " + path)
	}
	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalln("config file not parsable, path: " + path)
	}
	return cfg
}

func LoadBootstrapConfig() *Config {
	cfg := ReadConfig("config/bootstrap.yml")
	return &cfg
}
