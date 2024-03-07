package config

type Kafka struct {
	Addr            string `env:"MARS_BOOTSTRAP_SERVER"`
	Topic           string `env:"MARS_KAFKA_TOPIC"`
	MaxRetry        int    `env:"MARS_MAX_RETRY"`
	ReturnSuccesses bool   `env:"MARS_RETURN_SUCCESSES"`
}
