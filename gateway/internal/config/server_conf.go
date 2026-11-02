package config

type ServerConf struct {
	Addr string `yaml:"addr"`
	Mode int    `yaml:"mode"`
}
