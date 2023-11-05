package vv104

type Config struct {
	Mode     string `json:"mode"`
	Ipv4Addr string `json:"ipv4Addr"`
	Port     int    `json:"port"`
	Casdu    int    `json:"casdu"`
	AutoAck  bool   `json:"autoAck"`
	K        int    `json:"k"`
	W        int    `json:"w"`
	//t0 time.Time
	T1            int  `json:"t1"`
	T2            int  `json:"t2"`
	T3            int  `json:"t3"`
	IoaStructured bool `json:"ioaStructured"`
	UseLocalTime  bool `json:"useLocalTime"`
}

// const (
// 	server = iota
// 	client
// )
