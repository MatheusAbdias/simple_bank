package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Driver              string        `mapstructure:"DRIVER"`
	Source              string        `mapstructure:"DATABASE_URL"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	os.Setenv("DATABASE_URL", config.Source)

	return
}
