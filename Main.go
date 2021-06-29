package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	var counter = 0
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("tr", func(e *colly.HTMLElement) {

		if isUnitHeader(e) {
			counter++

			log.Println(e.ChildText("td"))
		}

		// e.Request.Visit(e.Attr("href"))
	})
	c.Visit("https://www.auravita.pl/mieszkania")
	log.Printf("ilosc mieszkan %v", counter)
}

func isUnitHeader(e *colly.HTMLElement) bool {
	classes := strings.Split(e.Attr("class"), " ")
	if len(classes) != 0 {
		for _, class := range classes {
			if class == "unit" {
				return true
			}
		}
	}
	return false
}
