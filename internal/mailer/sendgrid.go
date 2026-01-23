package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail string, apiKey string) *SendGridMailer {
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    sendgrid.NewSendClient(apiKey),
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, templateFile, body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})
	for i := 0; i < maxRetries; i++ {
		response, err := m.client.Send(message)
		if err != nil || (200 < response.StatusCode || response.StatusCode > 300) {
			if err != nil {
				log.Printf("Error %v %s", err, err.Error())
			}
			log.Printf("failed to send email, error %s, mail: %v attempt %d of %d status code %d", err, email, i+1, maxRetries, response.StatusCode)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("email sent successfully, attempt %d of %d response: %d", i+1, maxRetries, response.StatusCode)
		return nil
	}
	return fmt.Errorf("failed to send email, retries exceeded")
}
