package main

import (
	"log"

	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/roadmap"
)

func main() {
	overpassData, err := overpass.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}

	overpass.SummarizeOverpassData(overpassData)
	roadmap.TestSimplify()
}
