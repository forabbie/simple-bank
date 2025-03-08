package util

import (
	"fmt"
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

	viper.AutomaticEnv()

	// err = viper.ReadInConfig()
	// if err != nil {
	// 	return
	// }

	// Try to read the .env file, but don't fail if it's missing
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("No .env file found, using only environment variables")
		err = nil // Reset error so it does not stop execution
	}

	err = viper.Unmarshal(&config)
	return
}

// func LoadConfig(path string) (config Config, err error) {
// 	viper.AutomaticEnv() // Use env variables first

// 	if _, statErr := os.Stat(".env.test"); statErr == nil {
// 		viper.SetConfigFile(".env.test")
// 		err = viper.ReadInConfig() // Load from .env.test if it exists
// 		if err != nil {
// 			return
// 		}
// 	}

// 	err = viper.Unmarshal(&config)
// 	return
// }
