package overpass_import

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

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
