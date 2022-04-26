package zdpgo_env

type Config struct {
	Debug          bool   `json:"debug" yaml:"debug" env:"debug"`
	ConfigFilePath string `json:"config_file_path" yaml:"config_file_path" env:"config_file_path"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() (*Config, error) {
	return &Config{
		Debug:          true,
		ConfigFilePath: ".env",
	}, nil
}

// UpdateDefaultConfig 更新默认配置
func UpdateDefaultConfig(config Config) (cfg *Config) {
	if config.ConfigFilePath == "" {
		config.ConfigFilePath = ".env"
	}
	cfg = &config
	return
}
