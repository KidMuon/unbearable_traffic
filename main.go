package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/bearablemap"
	"github.com/KidMuon/unbearable_traffic/convert"
	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/resident"
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
	fmt.Println(streetmap["1738652113"])
	structuremap := convert.CreateBuildingMap(overpassData.BuildingResponse, streetmap)
	fmt.Println(len(structuremap))
	overpass.SummarizeOverpassData(overpassData)
	population := resident.CreatePopulation(9, structuremap)
	for _, pop := range population {
		fmt.Println(pop)
	}
}
