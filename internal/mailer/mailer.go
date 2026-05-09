package mailer

import (
	"bytes"
	"embed"
	"html/template"
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
		body:    subject.String(),
	}, nil
}
