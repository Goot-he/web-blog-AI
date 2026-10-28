package config

type LogConf struct {
	Level    string `yaml:"level"`
	Mode     string `yaml:"mode"`
	FileName string `yaml:"filename"`
}
