package overpass

import "fmt"

func SummarizeOverpassData(o OverpassData) {
	fmt.Printf("%d Streets in the area\n", o.StreetResponse.GetOverpassObjectCount())
	fmt.Printf("%d Buildings in the area\n", o.BuildingResponse.GetOverpassObjectCount())
}
