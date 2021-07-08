package stats

import (
	"github.com/janusz-chludzinski/aura-vita/models"
	"log"
	"strings"
	"sync"
)

const statusReserved = "rezerwacja"
const statusAvailable = "dostępny"

func GetStats(flats []models.Flat, data *models.MailData) *models.MailData {
	var wg sync.WaitGroup
	wg.Add(3)

	go count(flats, statusReserved, data, &wg)
	go count(flats, statusAvailable, data, &wg)
	go is5thFloorFull(flats, data, &wg)

	wg.Wait()
	return data
}

func is5thFloorFull(flats []models.Flat, data *models.MailData, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, flat := range flats {
		if flat.FlatNumber == "194" {
			data.AreAllNeighbours = false
		}
	}
}

func count(flats []models.Flat, status string, data *models.MailData, wg *sync.WaitGroup) {
	defer wg.Done()

	counter := 0
	for _, flat := range flats {
		if strings.ToLower(flat.Status) == status {
			counter++
		}
	}
	setData(data, status, counter)
}

func setData(data *models.MailData, status string, counter int) {
	switch status {
	case statusReserved:
		data.FlatsReserved = counter
	case statusAvailable:
		data.FlatsAvailable = counter
	default:
		log.Printf("Error: could not find match for status %v", status)
	}
}
