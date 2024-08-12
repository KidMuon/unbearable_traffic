package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/bearablemap"
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
	streetmap.Simplify()
	streetmap = bearablemap.EliminateDisconnectedNodes(streetmap)
	structuremap := convert.CreateBuildingMap(overpassData.BuildingResponse, streetmap)
	fmt.Println(len(structuremap))
	overpass.SummarizeOverpassData(overpassData)
	/*
		Create the people
		assign buildings to those people for work, leisure, home
		create schedules for the people to go to each location

	*/
}
