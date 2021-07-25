package models

import (
	"time"
)

type Flat struct {
	Building   string
	FlatNumber string
	Surface    string
	Status     string
	Floor      string
}

type ParseData struct {
	FlatsNotSold     int
	FlatsAvailable   int
	FlatsReserved    int
	AreAllNeighbours bool
	GalleriesCount   int
}

type DbEntry struct {
	ParseData *ParseData
	Date      time.Time
}
