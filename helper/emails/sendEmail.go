package emails

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/bishal-dd/receipt-generator-backend/pkg/rmq"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// LoadTemplate loads and executes an email template with provided data.
func LoadTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func SendEmail(to, subject, templatePath string, data interface{}) error {
	body, err := LoadTemplate(templatePath, data)
	if err != nil {
		return err
	}

	if rmq.EmailQueue == nil {
		return fmt.Errorf("email queue is not initialized")
	}

	email := EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	emailBytes, err := json.Marshal(email)
	if err != nil {
		return err
	}

	return rmq.EmailQueue.PublishBytes(emailBytes)
}
