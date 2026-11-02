package config

type RpcConf struct {
	CoreServerAddr string `yaml:"core_server_addr"`
	RequestTimeout int    `yaml:"request_timeout"`
}
