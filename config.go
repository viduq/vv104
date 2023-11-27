package vv104

import (
	"flag"
)

type Config struct {
	Mode            string `json:"mode"`
	Ipv4Addr        string `json:"ipv4Addr"`
	Port            int    `json:"port"`
	Casdu           int    `json:"casdu"`
	AutoAck         bool   `json:"autoAck"`
	K               int    `json:"k"`
	W               int    `json:"w"`
	T1              int    `json:"t1"`
	T2              int    `json:"t2"`
	T3              int    `json:"t3"`
	IoaStructured   bool   `json:"ioaStructured"`
	InteractiveMode bool   `json:"interactiveMode"`
	UseLocalTime    bool   `json:"useLocalTime"`
}

func (config *Config) ParseFlags() {

	serverPtr := flag.Bool("s", false, "Connection mode: For Server (Controlled station) use '-s'. For Client (Controlling station) use without flag (default)")
	ipPtr := flag.String("h", "127.0.0.1", "IP address")
	interactivePtr := flag.Bool("i", true, "Start in interactive mode, control program with cli commands")
	portPtr := flag.Int("p", 2404, "Port")
	flag.Parse()

	if *serverPtr {
		config.Mode = "server"
	} else {
		config.Mode = "client"
	}

	config.Ipv4Addr = *ipPtr
	config.InteractiveMode = *interactivePtr
	config.Port = *portPtr

	// todo
	config.AutoAck = true
	config.K = 12
	config.W = 8
	config.T1 = 15
	config.T2 = 10
	config.T3 = 20
	config.IoaStructured = false
	config.UseLocalTime = false

}
