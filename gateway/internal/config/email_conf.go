package config

type EmailConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Nickname string `yaml:"nickname"`
	Secret   string `yaml:"secret"`
	From     string `yaml:"from"`
	IsSSL    bool   `yaml:"is_ssl"`
}
