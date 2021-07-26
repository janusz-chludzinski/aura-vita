package scrapper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/janusz-chludzinski/aura-vita/models"
	"log"
	"net/http"
)

func GetFlats(url string) []models.Flat {
	doc, err := getHtml(url)
	if err != nil {
		log.Fatalf("[ERROR] while reading response: %v", err)
	}

	return getFlats(doc)
}

func GetPicsCount(url string) int {
	doc, err := getHtml(url)
	if err != nil {
		log.Fatalf("[ERROR] while reading response: %v", err)
	}

	return getPicsCount(doc)
}

func getPicsCount(doc *goquery.Document) int {
	return doc.Find(".gallery").Length()
}

func getHtml(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("[ERROR] while getting %v", url)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("[ERROR] Status code different then expected(200): %v", res.StatusCode)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func getFlats(doc *goquery.Document) []models.Flat {
	flats := make([]models.Flat, 0)

	log.Print("Collecting flats...")
	doc.Find(".unit").Each(func(i int, s *goquery.Selection) {
		flat := &models.Flat{}
		s.ChildrenFiltered("td").Each(func(i int, s *goquery.Selection) {
			value, _ := s.Attr("class")
			switch value {
			case "field-building":
				flat.Building = s.Text()
			case "title":
				flat.FlatNumber = s.Text()
			case "field-surface":
				flat.Surface = s.Text()
			case "field-availability":
				flat.Status = s.Text()
			case "field-floor":
				flat.Floor = s.Text()
			default:
				log.Printf("[INFO] No match found for %v, ignoring...", value)
			}
		})
		flats = append(flats, *flat)
	})
	log.Printf("[INFO] %v flats found!", len(flats))
	return flats
}
