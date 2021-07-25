package main

import (
	"context"
	"github.com/janusz-chludzinski/aura-vita/db"
	"github.com/janusz-chludzinski/aura-vita/mail"
	"github.com/janusz-chludzinski/aura-vita/models"
	"github.com/janusz-chludzinski/aura-vita/scrapper"
	"github.com/janusz-chludzinski/aura-vita/stats"
	. "go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/smtp"
	"time"
)

const flatsUrl = "https://www.auravita.pl/mieszkania"
const picsUrl = "https://www.auravita.pl/galeria"
const templatePath = "mail/template/email.html"
const dbName = "aura-vita"
const connectionString = "mongodb://admin:Start123@localhost:27017/aura-vita"

func main() {
	log.Println("[INFO] Starting task...")

	flats, count := scrappData(flatsUrl, picsUrl)
	data := getParseData(flats, count)
	_ = storeParseData(data)
	sendEmailNotification(data, mail.MailConfig{}.NewConfig())

	log.Println("[INFO] Task finished!")
}

func getAuth(config *mail.MailConfig) smtp.Auth {
	return smtp.PlainAuth("", config.From, config.Password, config.Host)
}

func getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func scrappData(flatsUrl string, picsUrl string) ([]models.Flat, int) {
	return scrapper.GetFlats(flatsUrl), scrapper.GetPicsCount(picsUrl)
}

func getParseData(flats []models.Flat, count int) *models.ParseData {
	data := &models.ParseData{}
	data.FlatsNotSold = len(flats)
	data.GalleriesCount = count
	stats.GetStats(flats, data)
	return data
}

func storeParseData(data *models.ParseData) *InsertOneResult {
	ctx := getContext()
	client, err := db.GetConnectedClient(connectionString, ctx)
	if err != nil {
		log.Println(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection("parse-result")
	entry := &models.DbEntry{
		ParseData: data,
		Date:      time.Now(),
	}
	savedEntry, err := db.Save(entry, collection, ctx)
	if err != nil {
		log.Println(err)
	}
	return savedEntry
}

func sendEmailNotification(data *models.ParseData, config *mail.MailConfig) {
	if err := mail.SendMail(templatePath, data, getAuth(config), config); err != nil {
		log.Printf("[ERROR]: could not send email. Reason: %v", err)
	}
}
