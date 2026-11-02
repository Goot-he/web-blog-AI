package config

type LogConf struct {
	Level    string `yaml:"level"`
	Mode     string `yaml:"mode"`
	Filename string `yaml:"filename"`
}
