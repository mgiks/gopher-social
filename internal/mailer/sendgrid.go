package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"time"

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

	letter, err := createLetter(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, letter.subject, to, "", letter.body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			// exponential backoff
			expBackoff := math.Pow(float64(time.Second), float64(i+i))
			time.Sleep(time.Duration(expBackoff))
			continue
		}

		log.Printf("Email sent with status code: %v\n", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
