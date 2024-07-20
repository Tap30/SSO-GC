package config

import (
	"github.com/spf13/viper"
)

// AppConfig holds the application configuration
type AppConfig struct {
	SsoIssuer string
}

// LoadConfig loads the application configuration using viper
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("sso_issuer", "https://accounts.backyard.tapsi.tech/api/v1/sso-user/oidc")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AppConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
