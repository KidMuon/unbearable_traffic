package data_import

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	overpassAPIURL              = "https://overpass-api.de/api/interpreter"
	boundingBoxExtentLimit      = 0.75
	overpassAPI_standard_config = "config/overpassAPI_standard_import.xml"
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
	Filters     []OverpassAPIFilter
	BoundingBox OverpassAPIBoundingBox
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
	for _, filter := range r.Filters {
		request += fmt.Sprintf(`["%s"~"%s"]`, filter.key, filter.value)
	}
	request += "\n"
	request += fmt.Sprintf(`(%f, %f, %f, %f);`, r.BoundingBox.south, r.BoundingBox.west, r.BoundingBox.north, r.BoundingBox.east)
	request += "\n(._;>;);\nout body;"
	return request
}

func ParseAsFilter(s string) (OverpassAPIFilter, error) {
	rawWords := strings.ReplaceAll(s, "[\"", "")
	rawWords = strings.ReplaceAll(rawWords, "\"~\"", " ")
	rawWords = strings.ReplaceAll(rawWords, "\"]", "")
	rawWords = strings.TrimSpace(rawWords)
	wordsSplit := strings.Split(rawWords, " ")

	if len(wordsSplit) != 2 {
		err := fmt.Sprintf("error parsing as filter. expected 2 words. got %d words", len(wordsSplit))
		return OverpassAPIFilter{}, errors.New(err)
	}

	filter := OverpassAPIFilter{
		key:   wordsSplit[0],
		value: wordsSplit[1],
	}

	return filter, nil
}

func ParseAsBoundingBox(s string) (OverpassAPIBoundingBox, error) {
	rawValues := strings.Trim(s, "()")
	listOfCoordinates := strings.Split(rawValues, ",")

	if len(listOfCoordinates) != 4 {
		return OverpassAPIBoundingBox{}, errors.New("incorrect number of values passed")
	}

	coordinateValues := make([]float32, 4)
	for i, coordinateString := range listOfCoordinates {
		coordinateCandidate, err := strconv.ParseFloat(strings.TrimSpace(coordinateString), 32)
		if err != nil {
			err_string := fmt.Sprintf("error converting %s to Float32: %s", coordinateString, err)
			return OverpassAPIBoundingBox{}, errors.New(err_string)
		}
		coordinateValues[i] = float32(coordinateCandidate)
	}

	return CreateBoundingBox(coordinateValues[0], coordinateValues[1], coordinateValues[2], coordinateValues[3])
}

func CreateBoundingBox(s, w, n, e float32) (OverpassAPIBoundingBox, error) {
	boundingBox := OverpassAPIBoundingBox{
		south: s,
		west:  w,
		north: n,
		east:  e,
	}

	if boundingBox.north-boundingBox.south > boundingBoxExtentLimit {
		return OverpassAPIBoundingBox{}, errors.New("north south extent is too far")
	}

	if boundingBox.east-boundingBox.west > boundingBoxExtentLimit {
		return OverpassAPIBoundingBox{}, errors.New("east west extent is too far")
	}

	if boundingBox.north < boundingBox.south {
		return OverpassAPIBoundingBox{}, errors.New("north boundary must be farther north than the south boundary")
	}

	if boundingBox.east < boundingBox.west {
		return OverpassAPIBoundingBox{}, errors.New("east boundary must be farther east than the west boundary")
	}

	return boundingBox, nil
}
