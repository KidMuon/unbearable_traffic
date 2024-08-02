package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/data_import"
)

func main() {
	/*
		nwr ["highway"~"."]["maxspeed"~"."]
		(33.64507, -112.37045, 33.65507, -112.36045);
		(._;>;);
		out body;
	*/

	region, err := data_import.CreateBoundingBox(33.64507, -112.37045, 33.65507, -112.36045)
	if err != nil {
		log.Fatal(err)
	}

	streetFilter1, err := data_import.ParseAsFilter(`["highway"~"."]`)
	if err != nil {
		log.Fatal(err)
	}
	streetFilter2, err := data_import.ParseAsFilter(`["maxspeed"~"."]`)
	if err != nil {
		log.Fatal(err)
	}
	streetRequest := data_import.OverpassAPIRequest{
		Filters:     []data_import.OverpassAPIFilter{streetFilter1, streetFilter2},
		BoundingBox: region,
	}
	streetResponse := data_import.OverpassStreetResponse{}
	streetResponse.BuildResponse(streetRequest)
	fmt.Println("Querying Streets...")
	fmt.Printf("%d Streets in the area\n", streetResponse.GetOverpassObjectCount())

	buildingFilter, err := data_import.ParseAsFilter(`["building"~"."]`)
	if err != nil {
		log.Fatal(err)
	}
	buildingRequest := data_import.OverpassAPIRequest{
		Filters:     []data_import.OverpassAPIFilter{buildingFilter},
		BoundingBox: region,
	}
	buildingResponse := data_import.OverpassBuildingResponse{}
	buildingResponse.BuildResponse(buildingRequest)
	fmt.Println("Querying Buildings...")
	fmt.Printf("%d Buildings in the area\n", buildingResponse.GetOverpassObjectCount())
}
