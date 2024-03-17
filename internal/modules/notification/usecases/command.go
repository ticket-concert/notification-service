package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"notification-service/configs"
	notification "notification-service/internal/modules/notification"
	"notification-service/internal/modules/notification/models/request"
	"notification-service/internal/modules/order"
	orderEntity "notification-service/internal/modules/order/models/entity"
	"notification-service/internal/pkg/errors"
	"notification-service/internal/pkg/helpers"
	"notification-service/internal/pkg/log"
	"os"
	"time"

	"github.com/skip2/go-qrcode"
	"go.elastic.co/apm"
)

var (
	Encode      = qrcode.Encode
	ParseFiles  = helpers.ParseFile
	GeneratePDF = helpers.GeneratePDF
	WriteFile   = os.WriteFile
	Remove      = os.Remove
)

type commandUsecase struct {
	logger               log.Logger
	helperMail           helpers.MailHelper
	orderRepositoryQuery order.MongodbRepositoryQuery
}

func NewCommandUsecase(log log.Logger, helperMail helpers.MailHelper, omq order.MongodbRepositoryQuery) notification.UsecaseCommand {
	return commandUsecase{
		logger:               log,
		helperMail:           helperMail,
		orderRepositoryQuery: omq,
	}
}

func (c commandUsecase) SendEmailOtpRegister(origCtx context.Context, payload request.SendEmailRegister) error {
	domain := "notificationUsecase-SendEmailOtpRegister"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	emailReq := helpers.MailRequest{
		From:    fmt.Sprintf("Soldev Concert Organizer<%s>", configs.GetConfig().Email.EmailUsername),
		To:      payload.Email,
		Subject: "Register User",
		Body:    fmt.Sprintf("Your otp register user is : %s", payload.Otp),
	}
	err := c.helperMail.SendEmail(emailReq)
	if err != nil {
		msg := "Error Send Email OTP Registration"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}
	return nil
}

func (c commandUsecase) SendEmailTicket(origCtx context.Context, orderId string) error {
	domain := "notificationUsecase-SendEmailTicket"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	orderData := <-c.orderRepositoryQuery.FindOrderById(ctx, orderId)
	if orderData.Error != nil {
		msg := "failed query order"
		c.logger.Error(ctx, orderData.Error.Error(), fmt.Sprintf("orderId: %+v", orderId))
		return errors.InternalServerError(msg)
	}

	if orderData.Data == nil {
		c.logger.Error(ctx, "order not found", fmt.Sprintf("orderId: %+v", orderId))
		return errors.BadRequest("order not found")
	}

	order, ok := orderData.Data.(*orderEntity.Order)
	if !ok {
		return errors.InternalServerError("failed marshal continentData")
	}

	emailRequest := request.EmailRequest{
		FullName:     order.FullName,
		EventName:    order.EventName,
		TicketNumber: order.TicketNumber,
		OrderTime:    order.OrderTime.Format("2006-01-02 15:04"),
		PaymentType:  order.Bank,
		EventTime:    order.DateTime.Format("2006-01-02 15:04"),
		EventPlace:   order.Country.Place,
	}

	qrCode, err := Encode(order.TicketNumber, qrcode.Medium, 170)
	if err != nil {
		msg := "failed encode qrcode"
		c.logger.Error(ctx, err.Error(), fmt.Sprintf("orderId: %+v", orderId))
		return errors.InternalServerError(msg)
	}
	fmt.Println(order.DateTime)
	ticketRequest := request.TicketRequest{
		FullName:     order.FullName,
		SeatNumber:   order.SeatNumber,
		EventName:    order.EventName,
		TicketType:   order.TicketType,
		TicketNumber: order.TicketNumber,
		TicketPrice:  fmt.Sprintf("$%d", order.Amount),
		EventTime:    order.DateTime.Format("2006-01-02 15:04"),
		EventPlace:   order.Country.Place,
		QRCode:       base64.StdEncoding.EncodeToString(qrCode),
	}

	emailTemplate := "template/notification.html"
	ticketTemplate := "template/ticket.html"
	body, err := ParseFiles(emailTemplate, emailRequest)
	if err != nil {
		msg := "Error Execute html"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(err.Error())
	}
	ticketPdf, err := GeneratePDF(ticketTemplate, ticketRequest)
	if err != nil {
		msg := "Error GeneratePDF"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}
	fileName := fmt.Sprintf("./temp/Ticket-Concert-Soldev_%s.pdf", order.TicketNumber)
	err = WriteFile(fileName, *ticketPdf, 0644)
	if err != nil {
		msg := "Error WriteFile"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	emailReq := helpers.MailRequest{
		From:    fmt.Sprintf("Soldev Concert Organizer<%s>", configs.GetConfig().Email.EmailUsername),
		To:      order.Email,
		Subject: fmt.Sprintf("Evoucher for %s", order.EventName),
		Body:    *body,
		Attach:  fileName,
	}
	err = c.helperMail.SendEmail(emailReq)
	if err != nil {
		msg := "Error Send Ticket Email"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	err = Remove(fileName)
	if err != nil {
		msg := fmt.Sprintf("Error Delete Attachment file: %s", order.TicketNumber)
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	c.logger.Info(ctx, fmt.Sprintf("Successful send email evoucher for %s", order.EventName), emailRequest)
	return nil
}
