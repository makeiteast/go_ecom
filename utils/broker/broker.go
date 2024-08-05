package broker

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Config holds all the configuration required for this package
type BrokerConfig struct {
	Host       string
	Port       string
	EmailTopic string
	SmsTopic   string
	Partition  int
}

// Returns new kafka Conn
func NewKafkaProducer(cfg *BrokerConfig) (smsConn *kafka.Conn, emailConn *kafka.Conn, err error) {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	smsConn, err = kafka.DialLeader(context.Background(), "tcp", address, cfg.SmsTopic, cfg.Partition)
	if err != nil {
		logrus.Error("failed to conn broker:", err)
		return nil, nil, err
	}

	emailConn, err = kafka.DialLeader(context.Background(), "tcp", address, cfg.EmailTopic, cfg.Partition)
	if err != nil {
		logrus.Fatal("failed to conn broker:", err)
		return nil, nil, err
	}

	return smsConn, emailConn, nil
}
