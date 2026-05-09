package mailer

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m SendGridMailer) Send(templateFile string, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	letter, err := constructLetter(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, letter.subject, to, "", letter.body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	return retry(maxRetries, func() (string, error) {
		_, err := m.client.Send(message)
		return email, err
	})
}
