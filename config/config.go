package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	Website              string        `mapstructure:"WEBSITE"`
	CrawlingTimeout      time.Duration `mapstructure:"CRAWLING_TIMEOUT"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RabbitMQAddress      string        `mapstructure:"RABBITMQ_ADDRESS"`
	ElasticSearchAddress string        `mapstructure:"ELASTICSEARCH_ADDRESS"`
	RunningInContainer   bool          `mapstructure:"RUNNING_IN_CONTAINER"`
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
