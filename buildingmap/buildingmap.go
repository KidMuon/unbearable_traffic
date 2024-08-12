package buildingmap

import (
	"math"

	"github.com/KidMuon/unbearable_traffic/roadmap"
)

type BuildingID string

type BuildingType string

var CommercialBuilding BuildingType = "commercial"
var LeisureBuilding BuildingType = "leisure"
var ResidentialBuilding BuildingType = "residential"
var SchoolBuilding BuildingType = "school"

type Building struct {
	Id              BuildingID
	ClosestRoadNode roadmap.NodeId
	BuildingType    BuildingType
	Nodes           []roadmap.Node
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

func (b *Building) AssignClosestRoadNode(r roadmap.RoadMap) {
	minDistance := 10.0
	thresholdDistance := 0.05
	var closest roadmap.NodeId

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
