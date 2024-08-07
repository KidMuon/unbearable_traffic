package overpass

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

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
