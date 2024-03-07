package kafka

import (
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	consumer         sarama.Consumer
	DocumentReceived chan struct{} // Канал для сигнала о получении документа
}

func NewConsumer(brokers []string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) ConsumeMessage(topic string) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to start consumer: ", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// Получаем сообщение из Kafka
			log.Printf("Received message: %s\n", string(msg.Value))

			// Обработка полученного сообщения

			// Отправляем сигнал о получении документа через канал
			c.DocumentReceived <- struct{}{}
			return
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
