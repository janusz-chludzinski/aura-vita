package main

import (
	"github.com/janusz-chludzinski/aura-vita/mail"
	"github.com/janusz-chludzinski/aura-vita/models"
	"log"
	"net/smtp"
)

const url = "https://www.auravita.pl/mieszkania"
const templatePath = "mail/template/email.html"

func main() {
	//scrapper.GetFlats(url)
	mailDataMock := models.MailData{
		FlatsAvailable: 5,
		FlatsReserved: 10,
		FlatsNotSold: 15,
		AreAllNeighbours: false,
	}
	config := mail.MailConfig{}.NewConfig()
	if err := mail.SendMail(templatePath, mailDataMock, getAuth(config), config); err != nil {
		log.Printf("Error: could not send email. Reason: %v", err)
	}
}

func getAuth(config *mail.MailConfig) smtp.Auth {
	return smtp.PlainAuth("", config.From, config.Password, config.Host)
}
