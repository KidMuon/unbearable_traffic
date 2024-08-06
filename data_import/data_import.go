package data_import

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
)

func TestPostOverpassAPI() error {
	body := `
		way ["highway"~"."]["maxspeed"~"."]
		(33.435, -112.08, 33.445, -112.07);
		(._;>;);
		out body;
	`

	resp, err := http.Post(overpassAPIURL, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	extractedWays := OverpassStreetResponse{}
	derr := xml.NewDecoder(resp.Body).Decode(&extractedWays)
	if derr != nil {
		return derr
	}

	fmt.Println(len(extractedWays.Streets))

	return nil
}

type Overpassapi_config struct {
	Bbox struct {
		South float32 `xml:"south"`
		West  float32 `xml:"west"`
		North float32 `xml:"north"`
		East  float32 `xml:"east"`
	} `xml:"bounding_box"`
	Street struct {
		Filters []struct {
			Key   string `xml:"k,attr"`
			Value string `xml:"v,attr"`
		} `xml:"filter"`
	} `xml:"street"`
	Building struct {
		Filters []struct {
			Key   string `xml:"k,attr"`
			Value string `xml:"v,attr"`
		} `xml:"filter"`
	} `xml:"building"`
}

func ImportOverpassDataUsingConfig() error {
	var overpass_config Overpassapi_config

	configXML, err := os.ReadFile("config/overpassAPI_standard_import.xml")
	if err != nil {
		return err
	}
	uerr := xml.Unmarshal(configXML, &overpass_config)
	if uerr != nil {
		return uerr
	}

	/*
		nwr ["highway"~"."]["maxspeed"~"."]
		(33.64507, -112.37045, 33.65507, -112.36045);
		(._;>;);
		out body;
	*/

	bb := overpass_config.Bbox
	region, err := CreateBoundingBox(bb.South, bb.West, bb.North, bb.East)
	if err != nil {
		return err
	}

	streetRequest := OverpassAPIRequest{
		BoundingBox: region,
	}
	for _, filter := range overpass_config.Street.Filters {
		filter_text := fmt.Sprintf(`["%s"~"%s"]`, filter.Key, filter.Value)
		streetFilter, err := ParseAsFilter(filter_text)
		if err != nil {
			return err
		}
		streetRequest.Filters = append(streetRequest.Filters, streetFilter)
	}
	streetResponse := OverpassStreetResponse{}
	streetResponse.BuildResponse(streetRequest)
	fmt.Println("Querying Streets...")
	fmt.Printf("%d Streets in the area\n", streetResponse.GetOverpassObjectCount())

	buildingRequest := OverpassAPIRequest{
		BoundingBox: region,
	}
	for _, filter := range overpass_config.Building.Filters {
		filter_text := fmt.Sprintf(`["%s"~"%s"]`, filter.Key, filter.Value)
		buildingFilter, err := ParseAsFilter(filter_text)
		if err != nil {
			return err
		}
		buildingRequest.Filters = append(buildingRequest.Filters, buildingFilter)
	}
	buildingResponse := OverpassBuildingResponse{}
	buildingResponse.BuildResponse(buildingRequest)
	fmt.Println("Querying Buildings...")
	fmt.Printf("%d Buildings in the area\n", buildingResponse.GetOverpassObjectCount())

	return nil
}
