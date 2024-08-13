package bearablemap

import (
	"math"
)

type BuildingID string

type BuildingType string

const (
	CommercialBuilding  BuildingType = "commercial"
	LeisureBuilding     BuildingType = "leisure"
	ResidentialBuilding BuildingType = "residential"
	SchoolBuilding      BuildingType = "school"
)

type Building struct {
	Id              BuildingID
	ClosestRoadNode NodeId
	BuildingType    BuildingType
	Nodes           []Node
	AverageLocation struct {
		Longitude float32
		Latitude  float32
	}
}

type BuildingMap map[BuildingID]Building

func (b *Building) AssignAverageLocation() {
	var sumLongitude float32 = 0.0
	var sumLatitude float32 = 0.0
	for _, node := range b.Nodes {
		sumLongitude += node.Longitude
		sumLatitude += node.Latitude
	}
	b.AverageLocation.Latitude = sumLatitude / float32(len(b.Nodes))
	b.AverageLocation.Longitude = sumLongitude / float32(len(b.Nodes))
}

func (b *Building) AssignClosestRoadNode(r RoadMap) {
	minDistance := 10.0
	thresholdDistance := 0.005 //Very Roughly half a kilometer
	var closest NodeId

	var distance float64
	for _, v := range r {
		a := b.AverageLocation.Longitude - v.Node.Longitude
		b := b.AverageLocation.Latitude - v.Node.Latitude
		distance = math.Hypot(float64(a), float64(b))
		if distance < minDistance {
			minDistance = distance
			closest = v.Node.Id
		}
	}

	if distance < thresholdDistance {
		b.ClosestRoadNode = closest
	}
}
