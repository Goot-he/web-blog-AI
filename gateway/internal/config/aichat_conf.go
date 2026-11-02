package config

type AiChatConf struct {
	VectorCollection string `yaml:"vector_collection"`
	DefaultTopK      int    `yaml:"default_top_k"`
}
