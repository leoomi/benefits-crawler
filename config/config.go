package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress   string        `mapstructure:"SERVER_ADDRESS"`
	Website         string        `mapstructure:"WEBSITE"`
	CrawlingTimeout time.Duration `mapstructure:"CRAWLING_TIMEOUT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	viper.Unmarshal(&config)
	return &config, nil
}
