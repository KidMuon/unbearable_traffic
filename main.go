package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/convert"
	"github.com/KidMuon/unbearable_traffic/overpass"
)

func main() {
	overpassData, err := overpass.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}
	waymap := convert.CreateWayMap(overpassData.StreetResponse)
	streetmap := convert.CreateRoadMap(waymap)
	intersectionCount := 0
	for _, v := range streetmap {
		if len(v.Edges) >= 3 {
			intersectionCount++
		}
	}
	fmt.Println("\n\n", intersectionCount, "\n\n------------")
	streetmap.Simplify()
	intersectionCount = 0
	for _, v := range streetmap {
		if len(v.Edges) >= 3 {
			intersectionCount++
		}
	}
	fmt.Println("\n\n", intersectionCount, "\n\n------------")
	overpass.SummarizeOverpassData(overpassData)
}
