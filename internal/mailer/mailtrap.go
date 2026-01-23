package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

// ✅ sandbox inbox id можно сделать константой
const mailtrapSandboxInboxID = 1234567 // <- поставь свой Inbox ID

type MailtrapMailer struct {
	fromEmail string
	apiKey    string
	client    *http.Client
}

func NewMailtrapMailer(fromEmail string, apiKey string) *MailtrapMailer {
	return &MailtrapMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type mailtrapAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailtrapSendRequest struct {
	From    mailtrapAddress   `json:"from"`
	To      []mailtrapAddress `json:"to"`
	Subject string            `json:"subject"`
	Text    string            `json:"text,omitempty"`
	HTML    string            `json:"html,omitempty"`
}

func (m *MailtrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mailtrapAddress{Email: m.fromEmail, Name: FromName}
	to := mailtrapAddress{Email: email, Name: username}

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	reqBody := mailtrapSendRequest{
		From:    from,
		To:      []mailtrapAddress{to},
		Subject: subject.String(),
		HTML:    body.String(),
		// Text:  можно добавить текстовый фоллбек при желании
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// endpoint: sandbox vs prod
	url := "https://send.api.mailtrap.io/api/send"
	if isSandbox {
		if mailtrapSandboxInboxID <= 0 {
			return fmt.Errorf("mailtrapSandboxInboxID must be > 0 for sandbox mode")
		}
		url = fmt.Sprintf("https://sandbox.api.mailtrap.io/api/send/%d", mailtrapSandboxInboxID)
	}

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer 9a8680736d77865055b85185b8186717")
		req.Header.Set("Content-Type", "application/json")

		resp, err := m.client.Do(req)
		if err != nil {
			log.Printf("mailtrap send error: %v", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		b, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		// ✅ успех: любой 2xx
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("email sent successfully via mailtrap, attempt %d of %d status %d", i+1, maxRetries, resp.StatusCode)
			return nil
		}
		log.Printf("failed to send email via mailtrap, to=%s attempt %d of %d status %d body=%s",
			email, i+1, maxRetries, resp.StatusCode, string(b))

		time.Sleep(time.Second * time.Duration(i+1))
	}

	return fmt.Errorf("failed to send email, retries exceeded")
}
