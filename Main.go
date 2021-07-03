package main

import "github.com/janusz-chludzinski/aura-vita/scrapper"

const url = "https://www.auravita.pl/mieszkania"

func main() {
	scrapper.GetFlats(url)
}
