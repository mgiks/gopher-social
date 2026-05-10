package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"math"
	"time"
)

const (
	FromName            = "GopherSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile string, username, email string, data any, isSandbox bool) error
}

type letter struct {
	subject string
	body    string
}

func constructLetter(templateFile string, data any) (letter, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return letter{}, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return letter{}, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return letter{}, err
	}

	return letter{
		subject: subject.String(),
		body:    body.String(),
	}, nil
}

func retry(retryCount int, sendEmail func() (string, error)) error {
	for i := range retryCount {
		email, err := sendEmail()
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of of %d\n", email, i+1, maxRetries)
			log.Printf("Error: %v\n", err.Error())

			// exponential backoff
			secsToWait := math.Pow(float64(2), float64(i))
			time.Sleep(time.Second * time.Duration(secsToWait))
			continue
		}
		log.Println("Email sent succesfully")
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
