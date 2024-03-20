package handlers

import (
	"fmt"
	"notification-service/configs"
	"notification-service/internal/modules/notification"
	kafkaConfluent "notification-service/internal/pkg/kafka/confluent"
	"notification-service/internal/pkg/log"
)

// confluent

func InitNotificationEventConflHandler(nc notification.UsecaseCommand, log log.Logger) {
	topicKcr := "asm-concert-send-otp-user-registration"
	fmt.Println(configs.GetConfig().ServiceName)
	kcr, _ := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), log)
	kcr.SetHandler(NewInsertNotificationConsumer(nc, log))
	kcr.Subscribe(topicKcr)

	topicTicket := "asm-concert-send-email-pdf"
	kte, _ := kafkaConfluent.NewConsumer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), log)
	kte.SetHandler(NewInsertNotificationConsumer(nc, log))
	kte.Subscribe(topicTicket)

}
