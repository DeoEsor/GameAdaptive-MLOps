package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func ConnectProducer(ctx context.Context) (sarama.SyncProducer, error) {
	brokerUrl := strings.Split(GetValue(KafkaBrokers), ",")
	logrus.Infof("Connecting in brockers %v", brokerUrl)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokerUrl, config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kafka producer: %w", err)
	}

	logrus.Infof("Connection in brockers %v successful!", brokerUrl)
	return conn, err
}
