package broker

import (
	"context"
	"go-store/internal/entity"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
)

type KafkaConn struct {
	*kafka.Conn
}

func NewUserBroker(conn *kafka.Conn) entity.AuthBroker {
	return &KafkaConn{conn}
}

func (b *KafkaConn) SendEmail(ctx context.Context, dest string, message []byte) error {
	bLog := log.WithFields(log.Fields{"func": "broker.SendEmail"})
	err := b.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	timeNowUnix := time.Now().Unix()
	timeNowStr := strconv.FormatUint(uint64(timeNowUnix), 10)
	_, err = b.Conn.WriteMessages(
		kafka.Message{
			Key:   []byte("email-" + timeNowStr),
			Value: message,
		},
	)
	if err != nil {
		bLog.Warning("failed to write messages:", err)
	}

	if err := b.Conn.Close(); err != nil {
		bLog.Warning("failed to close writer:", err)
	}

	return nil
}
