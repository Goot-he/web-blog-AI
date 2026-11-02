package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const filepath = "conf/config.yaml"

type Config struct {
	AiChat  AiChatConf  `yaml:"ai_chat"`
	Auth    AuthConf    `yaml:"auth"`
	Email   EmailConf   `yaml:"email"`
	Http    HttpConf    `yaml:"http"`
	Llm     LlmConf     `yaml:"lms"`
	Log     LogConf     `yaml:"log"`
	Mysql   MysqlConf   `yaml:"mysql"`
	Redis   RedisConf   `yaml:"redis"`
	Rpc     RpcConf     `yaml:"rpc"`
	Server  ServerConf  `yaml:"server"`
	Storage StorageConf `yaml:"storage"`
}

func LoadConfig() (*Config, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
