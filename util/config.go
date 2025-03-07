package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Read from environment variables first
	viper.AutomaticEnv()

	// Try to read from .env file (optional)
	if err = viper.ReadInConfig(); err != nil {
		log.Println("No .env file found. Using environment variables instead.")
		err = nil // Ignore error if .env is missing
	}

	// Unmarshal into config struct
	err = viper.Unmarshal(&config)
	return
}
