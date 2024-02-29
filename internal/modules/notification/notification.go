package notification

import (
	"context"
	notificationRequest "notification-service/internal/modules/notification/models/request"
)

type UsecaseCommand interface {
	SendEmailOtpRegister(origCtx context.Context, payload notificationRequest.SendEmailRegister) error
	SendEmailTicket(origCtx context.Context, orderId string) error
}

type MongodbRepositoryQuery interface {
}
