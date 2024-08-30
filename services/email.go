package services

import (
	"log"

	"gopkg.in/gomail.v2"
)

// EmailService holds the SMTP server configuration
type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendEmail sends an email with the specified subject and body to the recipient
func (es *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(es.Host, es.Port, es.Username, es.Password)

	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Could not send email to %s: %v", to, err)
		return err
	}
	log.Printf("Email sent to %s successfully!", to)
	return nil
}
