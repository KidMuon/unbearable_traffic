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
	streetmap := convert.CreateRoadMap(waymap)
	streetmap.Simplify()
	streetmap = bearablemap.EliminateDisconnectedNodes(streetmap)
	fmt.Println("Streets Mapped...")
	fmt.Println(len(streetmap))
	/* fmt.Println(streetmap["1738652113"]) */
	structuremap := convert.CreateBuildingMap(overpassData.BuildingResponse, streetmap)
	/* fmt.Println(len(structuremap)) */
	fmt.Println("Generating Population....")
	population := resident.CreatePopulation(1000, structuremap)
	fmt.Println("Population created. Generating Routes...")
	for _, pop := range population {
		pop.CalculateRoute(streetmap, 1)
	}
}
