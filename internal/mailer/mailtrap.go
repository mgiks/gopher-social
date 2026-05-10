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

func (m MailTrapMailer) NewSender(templateFile string, username, email string, data any, isSandbox bool) (Sender, error) {
	if isSandbox {
		return nil, nil
	}

	letter, err := constructLetter(templateFile, data)
	if err != nil {
		return nil, err
	}

	message := gomail.NewMessage()

	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", letter.subject)
	message.AddAlternative("text/html", letter.body)

	return sender{
		sender:        gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey),
		message:       message,
		receiverEmail: email,
	}, nil
}

type sender struct {
	message       *gomail.Message
	sender        *gomail.Dialer
	receiverEmail string
}

func (s sender) Send() (string, error) {
	return s.receiverEmail, s.sender.DialAndSend(s.message)
}
