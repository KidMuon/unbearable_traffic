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
	overpass.SummarizeOverpassData(overpassData)

	waymap := convert.CreateWayMap(overpassData.StreetResponse)
	fmt.Println("Ways Created.")

	streetmap := convert.CreateRoadMap(waymap)
	fmt.Println("Streets Created.")
	streetmap.Simplify()
	fmt.Println("Streets Simplified.")
	streetmap = bearablemap.EliminateDisconnectedNodes(streetmap)
	fmt.Println("Streets Finished Mapping.")
	fmt.Println(len(streetmap))

	structuremap := convert.CreateBuildingMap(overpassData.BuildingResponse, streetmap)
	fmt.Println("Building Map Created.")
	fmt.Println(len(structuremap))

	fmt.Println("Generating Population.")
	population := resident.CreatePopulation(1000, structuremap)
	fmt.Println("Population created. Generating Routes...")
	for _, pop := range population {
		pop.CalculateRoute(streetmap, 1)
	}
}
