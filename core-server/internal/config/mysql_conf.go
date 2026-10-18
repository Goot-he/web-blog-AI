package config

import "fmt"

type MysqlConf struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DbName      string `yaml:"dbname"`
	Config      string `yaml:"config"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
}

func (c *MysqlConf) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.Username, c.Password, c.Host, c.Port, c.DbName, c.Config)
}
