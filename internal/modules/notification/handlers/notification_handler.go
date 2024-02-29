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

type insertNotificationHandler struct {
	notificationUsecaseCommand notification.UsecaseCommand
}

func NewInsertNotificationConsumer(nc notification.UsecaseCommand) kafkaPkgConfluent.ConsumerHandler {
	return &insertNotificationHandler{
		notificationUsecaseCommand: nc,
	}
}

func (i insertNotificationHandler) SendEmailOtpRegister(message *k.Message, topic string) {
	log.GetLogger().Info(context.Background(), string(message.Value), fmt.Sprintf("Topic: %v Partition: %v - Offset: %v", *message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset.String()))

	var msg notificationRequest.SendEmailRegister
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.GetLogger().Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	if err := i.notificationUsecaseCommand.SendEmailOtpRegister(context.Background(), msg); err != nil {
		log.GetLogger().Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	return
}

func (i insertNotificationHandler) SendEmailTicket(message *k.Message, topic string) {
	log.GetLogger().Info(context.Background(), string(message.Value), fmt.Sprintf("Topic: %v Partition: %v - Offset: %v", *message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset.String()))

	var msg string
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.GetLogger().Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	if err := i.notificationUsecaseCommand.SendEmailTicket(context.Background(), msg); err != nil {
		log.GetLogger().Error(context.Background(), err.Error(), string(message.Value))
		return
	}
	return
}
