package mailer

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type MailTrapMailer struct {
	apiKey    string
	fromEmail string
}

func NewMailTrap(apiKey, fromEmail string) (MailTrapMailer, error) {
	if apiKey == "" {
		return MailTrapMailer{}, fmt.Errorf("api key is required")
	}
	return MailTrapMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}, nil
}

func (m MailTrapMailer) Send(templateFile string, username, email string, data any, isSandbox bool) error {
	if isSandbox {
		return nil
	}

	letter, err := constructLetter(templateFile, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()

	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", letter.subject)
	message.AddAlternative("text/html", letter.body)

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	return retry(maxRetries, func() (string, error) {
		return m.fromEmail, dialer.DialAndSend(message)
	})
}
