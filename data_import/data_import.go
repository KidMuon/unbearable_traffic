package data_import

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	//"errors"
)

const (
	overpassAPIURL = "https://overpass-api.de/api/interpreter"
)

type OverpassResponse interface {
	GetOverpassObjectCount() int
	BuildResponse(OverpassAPIRequest) error
}

type OverpassAPIWayNode struct {
	Reference_id string `xml:"ref,attr"`
}

type OverpassAPIWayTag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

type OverpassAPINode struct {
	Id  string `xml:"id,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
}

type OverpassAPIRequest struct {
	filters     []OverpassAPIFilter
	boundingBox OverpassAPIBoundingBox
}

type OverpassAPIFilter struct {
	key   string
	value string
}

type OverpassAPIBoundingBox struct {
	south float32
	west  float32
	north float32
	east  float32
}

func (r OverpassAPIRequest) GetString() (request string) {
	request += "nwr "
	for _, filter := range r.filters {
		request += fmt.Sprintf(`["%s"~"%s"]`, filter.key, filter.value)
	}
	request += fmt.Sprintf(`\n(%f, %f, %f, %f);`, r.boundingBox.south, r.boundingBox.west, r.boundingBox.north, r.boundingBox.east)
	request += "\n(._;>;);\nout body;"
	return request
}

type OverpassStreetResponse struct {
	Streets []OverpassAPIStreet `xml:"way"`
	Nodes   []OverpassAPINode   `xml:"node"`
}

type OverpassAPIStreet struct {
	Id          string               `xml:"id,attr"`
	StreetNodes []OverpassAPIWayNode `xml:"nd"`
	StreetTags  []OverpassAPIWayTag  `xml:"tag"`
}

func (o OverpassStreetResponse) GetOverpassObjectCount() int {
	return len(o.Streets)
}

func (o *OverpassStreetResponse) BuildResponse(request OverpassAPIRequest) error {
	body := request.GetString()

	resp, err := http.Post(overpassAPIURL, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	derr := xml.NewDecoder(resp.Body).Decode(o)
	if derr != nil {
		return derr
	}

	return nil
}

type OverpassBuildingResponse struct {
	Buildings []OverpassAPIBuilding `xml:"way"`
	Nodes     []OverpassAPINode     `xml:"node"`
}

type OverpassAPIBuilding struct {
	BuildingNodes []OverpassAPIWayNode `xml:"nd"`
	BuildingTags  []OverpassAPIWayTag  `xml:"tag"`
}

func (o OverpassBuildingResponse) GetOverpassObjectCount() int {
	return len(o.Buildings)
}

func (o *OverpassBuildingResponse) BuildResponse(request OverpassAPIRequest) error {
	body := request.GetString()

	resp, err := http.Post(overpassAPIURL, "text/plain;charset=UTF-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	derr := xml.NewDecoder(resp.Body).Decode(o)
	if derr != nil {
		return derr
	}

	return nil
}

// []APIfilters, boundingBox
// => OverpassResponses interface
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

	return nil
}
