package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/convert"
	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/roadmap"
)

func main() {
	overpassData, err := overpass.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(overpassData.StreetResponse.Streets[1])
	waymap := convert.CreateWayMap(overpassData.StreetResponse)
	fmt.Println(waymap["5590554"])
	overpass.SummarizeOverpassData(overpassData)
	roadmap.TestSimplify()
}
