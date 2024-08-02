package data_import

import "fmt"

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
