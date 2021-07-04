package mail

import (
	"bytes"
	"fmt"
	"github.com/janusz-chludzinski/aura-vita/models"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type MailConfig struct {
	Host string
	From string
	Password string
	Port string
	Mime string
	Subject string
	Receiver []string
}

func (MailConfig) NewConfig() *MailConfig {
	return &MailConfig{
		Host: "smtp.gmail.com",
		From: "elperro.gianni@gmail.com",
		Password: os.Getenv("GMAIL_PASS"),
		Port: "587",
		Mime: "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
		Subject: "Subject: " + "Test Email" + "!\n",
		Receiver: []string{"chludzinski.janusz@gmail.com"},
	}
}

func SendMail(template string, data models.MailData, auth smtp.Auth, config *MailConfig) error {
	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return err
	}
	return sendMail(config, auth, emailBody)
}

func parseTemplate(path string, data models.MailData) (string, error) {
	log.Printf("Parsing template %v", path)
	html, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error: could not process template from path %v", path)
	}

	buffer := new(bytes.Buffer)
	if err = html.Execute(buffer, data); err!= nil {
		return "", err
	}
	log.Printf("Template parsed!")
	return buffer.String(), nil
}

func getAddress(config *MailConfig) string{
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}

func getMessage(config *MailConfig, body string) []byte {
	return []byte(config.Subject + config.Mime + "\n" + body)
}

func sendMail(config *MailConfig, auth smtp.Auth, emailBody string) error{
	log.Print("Sending email...")
	err := smtp.SendMail(
		getAddress(config),
		auth,
		config.From,
		config.Receiver,
		getMessage(config, emailBody))
	if err != nil {
		return err
	}
	log.Print("Email sent!")
	return nil
}
