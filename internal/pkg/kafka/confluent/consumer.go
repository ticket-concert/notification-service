package kafka

import (
	"context"
	"fmt"
	"strings"

	"notification-service/internal/pkg/log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type consumer struct {
	handler  ConsumerHandler
	consumer *kafka.Consumer
	logger   log.Logger
}

// NewConsumer is a constructor of kafka consumer
func NewConsumer(cfg *kafka.ConfigMap, log log.Logger) (Consumer, error) {
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	return &consumer{
		logger:   log,
		consumer: c,
	}, nil
}

func (c *consumer) SetHandler(handler ConsumerHandler) {
	c.handler = handler
}

func (c *consumer) Subscribe(topics ...string) {
	if c.handler == nil {
		joinTopic := strings.Join(topics, ", ")
		msg := fmt.Sprintf("Kafka Consumer Error: Topics: [%s] There is no consumer handlers to handle message from incoming event", joinTopic)
		c.logger.Error(context.Background(), msg, fmt.Sprintf("%+v", topics))
		return
	}

	c.consumer.SubscribeTopics(topics, nil)
	go func() {
		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				msg := fmt.Sprintf("Kafka Consumer Error: %v (%v)\n", err, msg)
				c.logger.Error(context.Background(), msg, fmt.Sprintf("%+v", topics))
				continue
			}
			switch topics[0] {
			case "ticket-concert-send-otp-user-registration":
				go c.handler.SendEmailOtpRegister(msg, topics[0])
				c.consumer.CommitMessage(msg)
			case "ticket-concert-send-email-pdf":
				go c.handler.SendEmailTicket(msg, topics[0])
				c.consumer.CommitMessage(msg)
			default:
				c.consumer.CommitMessage(msg)
			}
		}
	}()

	return
}

func (c *consumer) Close(ctx context.Context) error {
	return c.consumer.Close()
}
