package handlers

import (
	"notification-service/internal/modules/notification"
	"notification-service/internal/modules/notification/models/request"
	"notification-service/internal/pkg/errors"
	"notification-service/internal/pkg/helpers"
	"notification-service/internal/pkg/log"
	"notification-service/internal/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

type NotificationHttpHandler struct {
	NotificationUsecaseCommand notification.UsecaseCommand
	Logger                     log.Logger
}

func InitNotificationHttpHandler(app *fiber.App, nuc notification.UsecaseCommand, log log.Logger, redisClient redis.Collections) {
	handler := &NotificationHttpHandler{
		NotificationUsecaseCommand: nuc,
		Logger:                     log,
	}
	route := app.Group("/api/notif")

	route.Get("/v1/send-ticket", handler.SendEmailTicket)
}

func (n NotificationHttpHandler) SendEmailTicket(c *fiber.Ctx) error {
	req := new(request.SendTicketReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, n.Logger, errors.BadRequest("bad request"))
	}

	err := n.NotificationUsecaseCommand.SendEmailTicket(c.Context(), req.OrderId)
	if err != nil {
		return helpers.RespCustomError(c, n.Logger, err)
	}
	return helpers.RespSuccess(c, n.Logger, "SendEmailTicket succes", "SendEmailTicket success")
}
