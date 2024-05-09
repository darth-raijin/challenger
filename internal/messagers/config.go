package messagers

type KafkaConfig struct {
	BootstrapServers []string `mapstructure:"KAFKA_BOOTSTRAP_SERVERS"`
	Topic            []string `mapstructure:"KAFKA_TOPICS"`
}
