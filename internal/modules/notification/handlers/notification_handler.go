package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"notification-service/internal/modules/notification"
	notificationRequest "notification-service/internal/modules/notification/models/request"
	kafkaPkgConfluent "notification-service/internal/pkg/kafka/confluent"
	"notification-service/internal/pkg/log"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	Unmarshal = json.Unmarshal
)

type NotificationEventHandler struct {
	NotificationUsecaseCommand notification.UsecaseCommand
	Logger                     log.Logger
}

func NewInsertNotificationConsumer(nc notification.UsecaseCommand, log log.Logger) kafkaPkgConfluent.ConsumerHandler {
	return &NotificationEventHandler{
		NotificationUsecaseCommand: nc,
		Logger:                     log,
	}
}

func (i NotificationEventHandler) SendEmailOtpRegister(message *k.Message, topic string) {
	i.Logger.Info(context.Background(), string(message.Value), fmt.Sprintf("Topic: %v Partition: %v - Offset: %v", *message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset.String()))

	var msg notificationRequest.SendEmailRegister
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		i.Logger.Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	if err := i.NotificationUsecaseCommand.SendEmailOtpRegister(context.Background(), msg); err != nil {
		i.Logger.Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	return
}

func (i NotificationEventHandler) SendEmailTicket(message *k.Message, topic string) {
	i.Logger.Info(context.Background(), string(message.Value), fmt.Sprintf("Topic: %v Partition: %v - Offset: %v", *message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset.String()))

	var msg string
	if err := Unmarshal(message.Value, &msg); err != nil {
		i.Logger.Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	if err := i.NotificationUsecaseCommand.SendEmailTicket(context.Background(), msg); err != nil {
		i.Logger.Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	return
}
