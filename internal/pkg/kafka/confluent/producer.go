package kafka

import (
	"context"
	"fmt"

	"notification-service/internal/pkg/log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer struct
type producer struct {
	producer *kafka.Producer
	logger   log.Logger
}

// NewProducer constructor
func NewProducer(cfg *kafka.ConfigMap, log log.Logger) (Producer, error) {
	p, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}

	prod := &producer{
		producer: p,
		logger:   log,
	}

	go prod.errReporter()

	return prod, nil
}

func (p *producer) errReporter() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				msg := fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition)
				p.logger.Error(context.Background(), msg, fmt.Sprintf("%+v", ev.TopicPartition.Error))
			}
		}
	}
}

func (p *producer) Publish(topic string, message []byte, kafkaPartition *int32) {
	partition := kafka.PartitionAny

	if kafkaPartition != nil {
		partition = *kafkaPartition
	}

	msgCh := p.producer.ProduceChannel()
	msgCh <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition,
		},
		Value: message,
	}
}

func (p *producer) Close(ctx context.Context) error {
	p.producer.Close()
	return nil
}
