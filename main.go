package main

import (
	"github.com/janusz-chludzinski/aura-vita/mail"
	"github.com/janusz-chludzinski/aura-vita/models"
	"github.com/janusz-chludzinski/aura-vita/scrapper"
	"github.com/janusz-chludzinski/aura-vita/stats"
	"log"
	"net/smtp"
)

const flatsUrl = "https://www.auravita.pl/mieszkania"
const picsUrl = "https://www.auravita.pl/galeria"
const templatePath = "mail/template/email.html"

func main() {
	flats := scrapper.GetFlats(flatsUrl)

	mailData := &models.MailData{}
	mailData.FlatsNotSold = len(flats)
	stats.GetStats(flats, mailData)
	mailData.GalleriesCount = scrapper.GetPicsCount(picsUrl)

	config := mail.MailConfig{}.NewConfig()
	if err := mail.SendMail(templatePath, mailData, getAuth(config), config); err != nil {
		log.Printf("Error: could not send email. Reason: %v", err)
	}
}

func getAuth(config *mail.MailConfig) smtp.Auth {
	return smtp.PlainAuth("", config.From, config.Password, config.Host)
}
