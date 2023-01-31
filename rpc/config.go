package rpc

type Config struct {
	Addr int    `json:"addr" yaml:"addr"`
	IP   string `yaml:"ip" json:"ip"`
	Name string `json:"name" yaml:"name"`
}
