package config

import (
	"github.com/labstack/gommon/log"

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

// LoadConfig initializes the AppConfig singleton instance with configuration data.
func LoadConfig() *AppConfig {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		instance = &AppConfig{
			// Default values
			SsoIssuer:    "YOUR_SSO_ISSUER_URL",
			ServerPort:   8080,
			ClientId:     "YOUR_CLIENT_ID",
			ClientSecret: "YOUR_CLIENT_SECRET",
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(instance); err != nil {
			log.Fatalf("Unable to decode into struct, %s", err)
		}
	})

	return instance
}
