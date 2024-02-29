package handlers

import (
	"fmt"
	"notification-service/configs"
	"notification-service/internal/modules/notification"
	kafkaConfluent "notification-service/internal/pkg/kafka/confluent"
	"notification-service/internal/pkg/log"
)

// confluent

func InitNotificationEventConflHandler(nc notification.UsecaseCommand) {
	topicKcr := "concert-send-otp-user-registration"
	fmt.Println(configs.GetConfig().ServiceName)
	kcr, _ := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), log.GetLogger())
	kcr.SetHandler(NewInsertNotificationConsumer(nc))
	kcr.Subscribe(topicKcr)

	topicTicket := "concert-send-email-pdf"
	kte, _ := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), log.GetLogger())
	kte.SetHandler(NewInsertNotificationConsumer(nc))
	kte.Subscribe(topicTicket)

}
