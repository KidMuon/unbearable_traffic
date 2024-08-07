package overpass_import

import (
	"fmt"
)

func ImportOverpassData(overpass_config OverpassAPI_config) (OverpassData, error) {

	bb := overpass_config.Bbox
	region, err := CreateBoundingBox(bb.South, bb.West, bb.North, bb.East)
	if err != nil {
		return OverpassData{}, err
	}

	streetRequest := OverpassAPIRequest{
		BoundingBox: region,
	}
	for _, filter := range overpass_config.Street.Filters {
		filter_text := fmt.Sprintf(`["%s"~"%s"]`, filter.Key, filter.Value)
		streetFilter, err := ParseAsFilter(filter_text)
		if err != nil {
			return OverpassData{}, err
		}
		streetRequest.Filters = append(streetRequest.Filters, streetFilter)
	}
	streetResponse := OverpassStreetResponse{}
	streetResponse.BuildResponse(streetRequest)

	buildingRequest := OverpassAPIRequest{
		BoundingBox: region,
	}
	for _, filter := range overpass_config.Building.Filters {
		filter_text := fmt.Sprintf(`["%s"~"%s"]`, filter.Key, filter.Value)
		buildingFilter, err := ParseAsFilter(filter_text)
		if err != nil {
			return OverpassData{}, err
		}
		buildingRequest.Filters = append(buildingRequest.Filters, buildingFilter)
	}
	buildingResponse := OverpassBuildingResponse{}
	buildingResponse.BuildResponse(buildingRequest)

	overpassData := OverpassData{
		StreetResponse:   streetResponse,
		BuildingResponse: buildingResponse,
	}

	return overpassData, nil
}

func ImportOverpassData_Standard() (OverpassData, error) {
	config, err := GetStandardConfig()
	if err != nil {
		return OverpassData{}, nil
	}
	return ImportOverpassData(config)
}
