package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	sarama.SyncProducer
}

func NewProducer(addrs []string) (*Producer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Retry.Backoff = 10 * time.Millisecond
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(addrs, cfg)
	if err != nil {
		return nil, fmt.Errorf("error init new kafka producer. %w", err)
	}

	return &Producer{
		SyncProducer: producer,
	}, nil
}

func (d *Producer) Publish(topic string, data []byte) error {
	_, _, err := d.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
