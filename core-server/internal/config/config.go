package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "conf/config.yaml"

type Config struct {
	Mysql  MysqlConf  `yaml:"mysql"`
	Redis  RedisConf  `yaml:"redis"`
	Log    LogConf    `yaml:"log"`
	Server ServerConf `yaml:"server"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}
	return cfg, nil
}
