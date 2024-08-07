package data_import

import (
	"encoding/xml"
	"os"
)

type OverpassAPI_config struct {
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

func GetStandardConfig() (OverpassAPI_config, error) {
	var overpass_config OverpassAPI_config

	configXML, err := os.ReadFile(overpassAPI_standard_config)
	if err != nil {
		return OverpassAPI_config{}, err
	}
	uerr := xml.Unmarshal(configXML, &overpass_config)
	if uerr != nil {
		return OverpassAPI_config{}, uerr
	}

	return overpass_config, nil
}
