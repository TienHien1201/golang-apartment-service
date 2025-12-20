package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Environment string

const (
	Dev  Environment = "dev"
	Qc   Environment = "qc"
	Prod Environment = "prod"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
	Data     DataConfig
	Ai       AiConfig
	JWT      JWTConfig
	Auth     AuthConfig
}

func LoadConfig(env Environment, configPath string) (*Config, error) {
	if env == "" {
		env = Dev // Default environment
	}

	// Load base config first
	baseConfig, err := loadBaseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	// Load environment specific config
	envConfig, err := loadEnvConfig(env, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s config: %w", env, err)
	}

	// Merge configurations
	config := mergeConfigs(baseConfig, envConfig)

	// Validate final configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func loadBaseConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("base")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	v.SetEnvPrefix("APARTMENT_APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func loadEnvConfig(env Environment, configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(string(env))
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	v.SetEnvPrefix("APARTMENT_APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func mergeConfigs(base, env *viper.Viper) *Config {
	var config Config

	// Merge base settings
	if err := base.Unmarshal(&config); err != nil {
		return nil
	}

	// Override with environment specific settings
	if err := env.Unmarshal(&config); err != nil {
		return nil
	}

	return &config
}

func validateConfig(cfg *Config) error {
	if cfg.Server.HTTP.Port <= 0 {
		return fmt.Errorf("invalid HTTP port")
	}

	return nil
}
