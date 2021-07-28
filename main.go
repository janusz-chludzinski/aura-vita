package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/janusz-chludzinski/aura-vita/db"
	"github.com/janusz-chludzinski/aura-vita/mail"
	"github.com/janusz-chludzinski/aura-vita/models"
	"github.com/janusz-chludzinski/aura-vita/scrapper"
	"github.com/janusz-chludzinski/aura-vita/stats"
	. "go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"time"
)

const flatsUrl = "https://www.auravita.pl/mieszkania"
const picsUrl = "https://www.auravita.pl/galeria"
const templatePath = "aura-vita/aura-vita/mail/template/email.html"
const dbName = "aura-vita"

func main() {
	log.SetOutput(os.Stdout)
	log.Println("[INFO] Starting task...")

	flats, count := scrappData(flatsUrl, picsUrl)
	data := getParseData(flats, count)
	parseDataToJsonFile(data)
	//_ = storeParseData(data)
	sendEmailNotification(data, mail.Config{}.NewConfig())

	log.Println("[INFO] Task finished!")
}

func getAuth(config *mail.Config) smtp.Auth {
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

func parseDataToJsonFile(data *models.ParseData) {
	currentTime := time.Now()

	entry := models.DbEntry{
		ParseData: data,
		Date:      currentTime,
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		log.Fatal("[ERROR] Could not marshal data.")
	}

	fileName := fmt.Sprintf("/home/pi/aura-vita/%v_parseEntry", currentTime.Format(time.RFC3339Nano))

	err = ioutil.WriteFile(fileName, entryJson, 0644)
	if err != nil {
		log.Printf("[ERROR] Could not store json entry: %v", err)
	}

	log.Printf("[INFO] Successfuly stored data: %s to file %v", entryJson, fileName)

}

func storeParseData(data *models.ParseData) *InsertOneResult {
	ctx := getContext()
	connStr, err := db.ConnectionString()

	if err != nil {
		log.Fatal(err)
	}

	client, err := db.GetConnectedClient(connStr, ctx)
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

func sendEmailNotification(data *models.ParseData, config *mail.Config) {
	//wd, _ := os.Getwd()
	if err := mail.SendMail(/*fmt.Sprintf("%v/%v", wd, templatePath)*/"mail/template/email.html", data, getAuth(config), config); err != nil {
		log.Printf("[ERROR]: could not send email. Reason: %v", err)
	}
}
