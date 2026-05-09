package mailer

import (
	"fmt"
	"log"
	"math"
	"time"

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

	letter, err := createLetter(templateFile, data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Subject: letter.subject,
		Html:    letter.body,
	}

	for i := range maxRetries {
		response, err := m.client.Emails.Send(params)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of of %d\n", email, i+1, maxRetries)
			log.Printf("Error: %v\n", err.Error())

			// exponential backoff
			secsToWait := math.Pow(float64(2), float64(i+1))
			time.Sleep(time.Second * time.Duration(secsToWait))
			continue
		}
		log.Printf("Email with id %v successfully sent!", response.Id)
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
