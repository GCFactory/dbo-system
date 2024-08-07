package config

type KafkaConsumer struct {
	Brokers string   `yaml:"brokers"`
	GroupID string   `yaml:"groupID"`
	Topics  []string `yaml:"topics"`
}

type KafkaProducer struct {
	Brokers string   `yaml:"brokers"`
	Topics  []string `yaml:"topics"`
}
