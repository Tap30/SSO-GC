package config

import (
	"github.com/spf13/viper"
	"sync"
)

type AppConfig struct {
	SsoIssuer    string
	ServerPort   int    `mapstructure:"server_port"`
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

var (
	instance *AppConfig
	once     sync.Once
)

func LoadConfig() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{
			SsoIssuer:    "YOUR_SSO_ISSUER_URL",
			ServerPort:   8080,
			ClientId:     "YOUR_CLIENT_ID",
			ClientSecret: "YOUR_CLIENT_SECRET",
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
