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
	//fmt.Println(waymap["5590554"])
	streetmap := convert.CreateRoadMap(waymap)
	//fmt.Println(streetmap["41524705"])
	fmt.Println(len(overpassData.StreetResponse.Nodes))
	fmt.Println(len(streetmap))
	streetmap.Simplify()
	fmt.Println(len(streetmap))
	overpass.SummarizeOverpassData(overpassData)
	roadmap.TestSimplify()
}
