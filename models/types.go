package models

type Flat struct {
	Building   string
	FlatNumber string
	Surface    string
	Status     string
	Floor      string
}

type MailData struct {
	FlatsNotSold     int
	FlatsAvailable   int
	FlatsReserved    int
	AreAllNeighbours bool
}
