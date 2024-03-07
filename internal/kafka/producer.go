package kafka

import (
	"github.com/IBM/sarama"
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/config"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(cfg *config.Kafka) (*Producer, error) {
	KafkaCfg := cfg

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = KafkaCfg.MaxRetry
	config.Producer.Return.Successes = KafkaCfg.ReturnSuccesses

	producer, err := sarama.NewSyncProducer([]string{KafkaCfg.Addr}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) SendReport(topic string, message string) error {
	report := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := p.producer.SendMessage(report)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
