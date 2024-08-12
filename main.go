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
	waymap := convert.CreateWayMap(overpassData.StreetResponse)
	streetmap := convert.CreateRoadMap(waymap)
	streetmap.Simplify()
	streetmap = roadmap.EliminateDisconnectedNodes(streetmap)
	fmt.Println(len(streetmap))
	fmt.Println(streetmap)
	overpass.SummarizeOverpassData(overpassData)
}
