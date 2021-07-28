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

type Config struct {
	Host     string
	From     string
	Password string
	Port     string
	Mime     string
	Subject  string
	Receiver []string
}

func (Config) NewConfig() *Config {
	return &Config{
		Host:     "smtp.gmail.com",
		From:     "elperro.gianni@gmail.com",
		Password: os.Getenv("GMAIL_PASS"),
		Port:     "587",
		Mime:     "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
		Subject:  "Subject: " + "Aura-Vita status update!",
		Receiver: []string{"chludzinski.janusz@gmail.com"},
	}
}

func SendMail(template string, data *models.ParseData, auth smtp.Auth, config *Config) error {
	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return err
	}
	return sendMail(config, auth, emailBody)
}

func parseTemplate(path string, data *models.ParseData) (string, error) {
	log.Printf("[INFO] Parsing template %v", path)
	html, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("[ERROR] could not process template from path %v", path)
	}

	buffer := new(bytes.Buffer)
	if err = html.Execute(buffer, data); err != nil {
		return "", err
	}
	log.Printf("[INFO] Template parsed!")
	return buffer.String(), nil
}

func getAddress(config *Config) string {
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}

func getMessage(config *Config, body string) []byte {
	return []byte(config.Subject + config.Mime + "\n" + body)
}

func sendMail(config *Config, auth smtp.Auth, emailBody string) error {
	log.Print("[INFO] Sending email...")
	err := smtp.SendMail(
		getAddress(config),
		auth,
		config.From,
		config.Receiver,
		getMessage(config, emailBody))
	if err != nil {
		return err
	}
	log.Print("[INFO] Email sent!")
	return nil
}
