package helpers

import (
	"fmt"
	"net/smtp"
)

type EmailImpl struct{}

// EmailSender is an interface for sending emails.
type EmailSender interface {
	SendEmailSimple(sbjt string, to string, msg string) error
	Init(username string, password string)
}

var (
	emailSender EmailSender
	userEmail   string
	passEmail   string
)

func (e *EmailImpl) Init(username string, password string) {
	userEmail = username
	passEmail = password
}

func (e *EmailImpl) SendEmailSimple(sbjt string, to string, msg string) error {
	auth := smtp.PlainAuth(
		"",
		userEmail,
		passEmail,
		"smtp.gmail.com",
	)

	messages := fmt.Sprintf("Subject: %s\n\n%s", sbjt, msg)
	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		userEmail,
		[]string{to},
		[]byte(messages))
}
