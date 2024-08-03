package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
	"strconv"
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
		viper.AddConfigPath("./config/")

		instance = &AppConfig{
			// Default values
			SsoIssuer:    getEnv("SSO_ISSUER", "YOUR_SSO_ISSUER_URL"),
			ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
			ClientId:     getEnv("CLIENT_ID", "YOUR_CLIENT_ID"),
			ClientSecret: getEnv("CLIENT_SECRET", "YOUR_CLIENT_SECRET"),
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("Error reading config file, %s", err)
		} else {
			if err := viper.Unmarshal(instance); err != nil {
				log.Fatalf("Unable to decode into struct, %s", err)
			}
		}
	})

	return instance
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultVal
}
