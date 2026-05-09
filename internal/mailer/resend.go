package mailer

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	apiKey    string
	fromEmail string
	client    *resend.Client
}

func NewResend(apiKey, fromEmail string) (ResendMailer, error) {
	if apiKey == "" {
		return ResendMailer{}, fmt.Errorf("api key is required")
	}
	client := resend.NewClient(apiKey)
	return ResendMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    client,
	}, nil
}

func (m ResendMailer) Send(templateFile string, username, email string, data any, isSandbox bool) error {
	if isSandbox {
		return nil
	}

	letter, err := constructLetter(templateFile, data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Subject: letter.subject,
		Html:    letter.body,
	}

	return retry(maxRetries, func() (string, error) {
		_, err := m.client.Emails.Send(params)
		return email, err
	})
}
