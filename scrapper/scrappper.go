package scrapper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetFlats(url string) []Flat {
	doc, err := getHtml(url)

	if err != nil {
		log.Fatalf("Error while reading response: %v", err)
	}
	return getFlats(doc)
}

func getHtml(url string) (*goquery.Document, error) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error while getting %v", url)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error! Status code different then expected(200): %v", res.StatusCode)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func getFlats(doc *goquery.Document) []Flat {
	flats := make([]Flat, 0)

	log.Print("Collecting flats...")
	doc.Find(".unit").Each(func(i int, s *goquery.Selection) {
		flat := &Flat{}
		s.ChildrenFiltered("td").Each(func(i int, s *goquery.Selection) {
			value, _ := s.Attr("class")

			switch value {
			case "field-building":
				flat.building = s.Text()
			case "title":
				flat.flatNumber = s.Text()
			case "field-surface":
				flat.surface = s.Text()
			case "field-availability":
				flat.status = s.Text()
			default:
				log.Printf("No match found for %v, ignoring...", value)
			}
			flats = append(flats, *flat)
		})
	})
	log.Printf("%v flats found!", len(flats))
	return flats
}
