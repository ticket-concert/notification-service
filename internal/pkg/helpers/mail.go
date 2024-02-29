package helpers

import (
	"fmt"
	"log"
	"notification-service/configs"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	Email      string
	Password   string
	SenderName string
	SmtpHost   string
	SmtpPort   string
}

type MailRequest struct {
	From    string
	To      []string
	CC      string
	CCName  string
	Subject string
	Body    string
	Attach  string
}

func (m *Mail) SendEmail(payload MailRequest) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fmt.Sprintf("Soldev Concert <%s>", payload.From))
	mailer.SetHeader("To", payload.To...)
	mailer.SetAddressHeader("Cc", payload.CC, payload.CCName)
	mailer.SetHeader("Subject", payload.Subject)
	mailer.SetBody("text/html", payload.Body)
	mailer.Attach(payload.Attach)

	fmt.Printf("Soldev Concert <%s>", payload.From)
	dialer := gomail.NewDialer(
		configs.GetConfig().Email.SmtpHost,
		587,
		configs.GetConfig().Email.EmailUsername,
		configs.GetConfig().Email.EmailPassword,
	)

	fmt.Println(configs.GetConfig().Email.SmtpHost)
	fmt.Println(configs.GetConfig().Email.SmtpPort)
	fmt.Println(configs.GetConfig().Email.EmailUsername)
	fmt.Println(configs.GetConfig().Email.EmailPassword)

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
