package vv104

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Mode            string `toml:"mode"`
	Ipv4Addr        string `toml:"ipv4Addr"`
	Port            int    `toml:"port"`
	Casdu           int    `toml:"casdu"`
	AutoAck         bool   `toml:"autoAck"`
	K               int    `toml:"k"`
	W               int    `toml:"w"`
	T1              int    `toml:"t1"`
	T2              int    `toml:"t2"`
	T3              int    `toml:"t3"`
	IoaStructured   bool   `toml:"ioaStructured"`
	InteractiveMode bool   `toml:"interactiveMode"`
	UseLocalTime    bool   `toml:"useLocalTime"`
	LogToBuffer     bool   `toml:"logToBuffer"`
	LogToStdOut     bool   `toml:"logToStdOut"`
}

type eximportConfig struct {
	Config            *Config
	ConfiguredObjects *ConfiguredObjects
}

func NewConfig() *Config {
	cfg := Config{}
	return &cfg
}

func (config *Config) ParseFlags(objects *Objects) {

	confPathPtr := flag.String("c", "", "Path to config toml file. If conf.toml is provided, all other flags are overwritten")
	serverPtr := flag.Bool("s", false, "Connection mode: For Server (Controlled station) use '-s'. For Client (Controlling station) use without flag (default)")
	ipPtr := flag.String("h", "127.0.0.1", "IP address")
	interactivePtr := flag.Bool("i", true, "Start in interactive mode, control program with cli commands")
	portPtr := flag.Int("p", 2404, "Port")
	logToBufferPtr := flag.Bool("b", false, "Log to buffer")
	logToStdOutPtr := flag.Bool("l", true, "Log to standard output")

	flag.Parse()

	if *confPathPtr != "" {
		// use conf.toml

		loadedConfig, loadedObjects, err := LoadConfigAndObjectsFromFile(*confPathPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		*config = *loadedConfig
		*objects = *loadedObjects

	} else {

		if *serverPtr {
			config.Mode = "server"
		} else {
			config.Mode = "client"
		}

		config.Ipv4Addr = *ipPtr
		config.InteractiveMode = *interactivePtr
		config.Port = *portPtr
		config.LogToBuffer = *logToBufferPtr
		config.LogToStdOut = *logToStdOutPtr

		// todo
		config.K = 12
		config.W = 8
		config.T1 = 15
		config.T2 = 10
		config.T3 = 20
		config.AutoAck = true
		config.IoaStructured = false
		config.UseLocalTime = false
	}

}

func printConfig(config Config) {
	logInfo.Println("============= Config =============")
	logInfo.Printf("%+v\n", config)
}

func WriteConfigAndObjectsToFile(config Config, objects Objects, filePathAndName string) error {

	eximportConfig := &eximportConfig{
		Config:            &config,
		ConfiguredObjects: &objects.configuredObjects,
	}

	data, err := toml.Marshal(eximportConfig)
	if err != nil {
		logError.Println(err)
		return err
	}

	fmt.Printf("%s\n", data)

	if filePathAndName[len(filePathAndName)-5:] != ".toml" {
		fmt.Println("Please save as .toml file")
		return errors.New("please save as .toml file")
	}

	err = os.WriteFile(filepath.Join(filePathAndName), data, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func LoadConfigAndObjectsFromFile(filePathAndName string) (*Config, *Objects, error) {

	eximportConfig := &eximportConfig{
		Config:            &Config{},
		ConfiguredObjects: &ConfiguredObjects{},
	}

	data, err := os.ReadFile(filePathAndName)
	if err != nil {
		return nil, nil, err
	}

	err = toml.Unmarshal(data, eximportConfig)
	if err != nil {
		return nil, nil, err
	}

	// fmt.Println(*eximportConfig.ConfiguredObjects)

	objects := NewObjects()

	err = objects.addObjectsFromList(*eximportConfig.ConfiguredObjects)
	if err != nil {
		return nil, nil, err
	}

	return eximportConfig.Config, objects, nil
}
