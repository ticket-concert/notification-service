package usecases

import (
	"bytes"
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
	"strconv"
	"text/template"
	"time"

	"github.com/skip2/go-qrcode"
	"go.elastic.co/apm"
	"gopkg.in/gomail.v2"
)

type commandUsecase struct {
	logger               log.Logger
	helperEmail          helpers.EmailSender
	helperMail           helpers.Sender
	orderRepositoryQuery order.MongodbRepositoryQuery
}

func NewCommandUsecase(log log.Logger, helperEmail helpers.EmailSender, helperMail helpers.Sender, omq order.MongodbRepositoryQuery) notification.UsecaseCommand {
	return commandUsecase{
		logger:               log,
		helperEmail:          helperEmail,
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
	err := c.helperEmail.SendEmailSimple("Register User", payload.Email, fmt.Sprintf("Your otp register user is : %s", payload.Otp))
	if err != nil {
		msg := "Send email failed"
		c.logger.Error(ctx, err.Error(), fmt.Sprintf("%+v", payload))
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

	qrCode, err := qrcode.Encode(order.TicketNumber, qrcode.Medium, 170)
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
	// emailMsg, err :=
	t, err := template.ParseFiles(emailTemplate)
	if err != nil {
		msg := "Error ParseFiles html"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, emailRequest); err != nil {
		msg := "Error Execute html"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}
	body := buf.String()
	ticketPdf, err := helpers.GeneratePDF(ticketTemplate, ticketRequest)
	if err != nil {
		msg := "Error GeneratePDF"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}
	fileName := fmt.Sprintf("./temp/Ticket-Concert-Soldev_%s.pdf", order.TicketNumber)
	err = os.WriteFile(fileName, *ticketPdf, 0644)
	if err != nil {
		msg := "Error WriteFile"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	// mailPayload := helpers.MailRequest{
	// 	From:    configs.GetConfig().Email.EmailUsername,
	// 	To:      []string{"alifsnitt@gmail.com"},
	// 	Subject: fmt.Sprintf("Evoucher for %s", order.EventName),
	// 	Body:    body,
	// 	Attach:  fileName,
	// }

	// fmt.Println(mailPayload.From)
	// fmt.Println(mailPayload.To)
	// emailSender := helpers.New(configs.GetConfig().Email.SmtpHost, configs.GetConfig().Email.EmailPassword, configs.GetConfig().Email.EmailUsername, configs.GetConfig().Email.EmailPassword)
	// fmt.Println(body)
	// subject := fmt.Sprintf("Evoucher for %s", order.EventName)
	// msg := helpers.NewMessage(subject, body)
	// msg.To = []string{"alifsnitt@gmail.com"}
	// msg.AttachFile(fileName)

	// err = c.helperMail.Send(msg)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("Soldev Concert Organizer<%s>", configs.GetConfig().Email.EmailUsername))
	mailer.SetHeader("To", order.Email)
	mailer.SetHeader("Subject", fmt.Sprintf("Evoucher for %s", order.EventName))
	mailer.SetBody("text/html", body)
	mailer.Attach(fileName)

	smtpPort, _ := strconv.Atoi(configs.GetConfig().Email.SmtpPort)
	dialer := gomail.NewDialer(
		configs.GetConfig().Email.SmtpHost,
		smtpPort,
		configs.GetConfig().Email.EmailUsername,
		configs.GetConfig().Email.EmailPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		msg := "Error Send Ticket Email"
		c.logger.Error(ctx, msg, err.Error())
		return errors.InternalServerError(msg)
	}

	return nil
}
