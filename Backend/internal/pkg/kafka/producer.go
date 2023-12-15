package kafka

import (
	"context"
	"encoding/json"
	"fmt"

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
func (k kafkaProduser) SendMessage(ctx context.Context, key string, message interface{}) error {
	messageToSend, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal message: %w", err)
	}
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(string(messageToSend)),
	}

	partition, offset, err := k.producerClient.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("cannot send message to kafka: %w", err)
	}

	logrus.Infof("Message is stored in topic=%s, partition=%d; offset is %d", k.topic, partition, offset)

	return nil
}
