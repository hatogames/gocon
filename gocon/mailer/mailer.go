package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// SMTPConfig enthält die SMTP-Server Einstellungen
var (
	smtpHost     = "smtp.mail.yahoo.com" // z.B. "smtp.gmail.com"
	smtpPort     = 465                   // 465 = SSL, 587 = STARTTLS
	smtpUsername = "hatogames@yahoo.com"
	smtpPassword = "peicgejjhhpiuurg"
	smtpSSL      = true
	fromAddress  = smtpUsername
)

func SendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
	d.SSL = smtpSSL

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("Fehler beim Senden der E-Mail: %w", err)
	}

	return nil
}
