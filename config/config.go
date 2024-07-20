package config

import (
	"github.com/spf13/viper"
	"sync"
)

type AppConfig struct {
	SsoIssuer  string
	ServerPort int `mapstructure:"server_port"`
}

var (
	instance *AppConfig
	once     sync.Once
)

func LoadConfig() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{
			SsoIssuer:  "https://development.backyard.tapsi.tech/api/v2/user/sso",
			ServerPort: 8080,
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err == nil {
			viper.Unmarshal(instance)
		}
	})

	return instance
}
