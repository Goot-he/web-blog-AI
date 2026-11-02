package config

type StorageConf struct {
	Provider  string `yaml:"provider"`
	RootDir   string `yaml:"root_dir"`
	BaseUrl   string `yaml:"base_url"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	MaxSize   int    `yaml:"max_size"`
}
