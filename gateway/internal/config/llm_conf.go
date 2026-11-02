package config

type LlmConf struct {
	Provider  string `yaml:"provider"`
	ApiKey    string `yaml:"api_key"`
	ApiUrl    string `yaml:"api_url"`
	ModelId   string `yaml:"model_id"`
	ModelName string `yaml:"model_name"`
}
