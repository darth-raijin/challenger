package config

import (
	"github.com/darth-raijin/challenger/internal/messagers"
	"github.com/spf13/viper"
)

// Config stores configuration data
type Config struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	Kafka         messagers.KafkaConfig
	KafkaProducer messagers.ProducerConfig
	TestKafka     []string `mapstructure:"KAFKA_TOPICS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
