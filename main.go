package main

import (
	"log"

	"github.com/KidMuon/unbearable_traffic/overpass_import"
)

func main() {
	overpassData, err := overpass_import.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}
	overpass_import.SummarizeOverpassData(overpassData)

}
