package helpers

import (
	"log"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	Email      string
	Password   string
	SenderName string
	SmtpHost   string
	SmtpPort   int
}

type MailRequest struct {
	From    string
	To      string
	Subject string
	Body    string
	Attach  string
}

func (m *Mail) SendEmail(payload MailRequest) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", payload.From)
	mailer.SetHeader("To", payload.To)
	mailer.SetHeader("Subject", payload.Subject)
	mailer.SetBody("text/html", payload.Body)
	if len(payload.Attach) > 0 {
		mailer.Attach(payload.Attach)
	}

	dialer := gomail.NewDialer(
		m.SmtpHost,
		m.SmtpPort,
		m.Email,
		m.Password,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

type MailHelper interface {
	SendEmail(payload MailRequest) error
}
