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

const BuildingThresholdDistance float64 = 0.005 //Roughly half a kilometer

type Building struct {
	Id              BuildingID
	ClosestRoadNode NodeId
	BuildingType    BuildingType
	Nodes           []Node
	AverageLocation struct {
		Longitude float64
		Latitude  float64
	}
}

type BuildingMap map[BuildingID]Building

func (b *Building) AssignAverageLocation() {
	var sumLongitude float64 = 0.0
	var sumLatitude float64 = 0.0
	for _, node := range b.Nodes {
		sumLongitude += node.Longitude
		sumLatitude += node.Latitude
	}
	b.AverageLocation.Latitude = sumLatitude / float64(len(b.Nodes))
	b.AverageLocation.Longitude = sumLongitude / float64(len(b.Nodes))
}

func (b *Building) AssignClosestRoadNode(r RoadMap) {
	minDistance := 10.0
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

	if minDistance < BuildingThresholdDistance {
		b.ClosestRoadNode = closest
	}
}
