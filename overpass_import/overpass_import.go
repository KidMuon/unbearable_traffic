package overpass_import

import (
	"fmt"
)

func ImportOverpassData(overpass_config OverpassAPI_config) error {

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

func ImportOverpassData_Standard() error {
	config, err := GetStandardConfig()
	if err != nil {
		return nil
	}
	return ImportOverpassData(config)
}
