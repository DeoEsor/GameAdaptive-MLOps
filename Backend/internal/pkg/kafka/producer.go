package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaProduser struct {
	topic          string
	producerClient sarama.SyncProducer
}

func NewKafkaProducer(
	topic string,
	producerClient sarama.SyncProducer,
) kafkaProduser {
	return kafkaProduser{
		topic:          topic,
		producerClient: producerClient,
	}
}

// SendMessage - отправляет сообщение в кафку
func (k kafkaProduser) SendMessage(ctx context.Context, message time.Time) error {
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(message.Format(time.RFC3339)),
	}

	partition, offset, err := k.producerClient.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("cannot send message to kafka: %w", err)
	}

	logrus.Infof("Message is stored in topic=%s, partition=%d; offset is %d", k.topic, partition, offset)

	return nil
}
