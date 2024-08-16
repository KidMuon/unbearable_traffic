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

func (b *Building) AssignClosestRoadNode(r RoadMap, roadindex RoadIndex) {
	minimumFoundDistance := 10.0
	var closestFoundNode NodeId

	searchPoints := []IndexPoint{
		{
			IndexLongitude: int(b.AverageLocation.Longitude * 100),
			IndexLatitude:  int(b.AverageLocation.Latitude * 100),
		},
		{
			IndexLongitude: int(b.AverageLocation.Longitude * 100),
			IndexLatitude:  int(b.AverageLocation.Latitude * 100),
		},
		{
			IndexLongitude: int(b.AverageLocation.Longitude * 100),
			IndexLatitude:  int(b.AverageLocation.Latitude * 100),
		},
		{
			IndexLongitude: int(b.AverageLocation.Longitude * 100),
			IndexLatitude:  int(b.AverageLocation.Latitude * 100),
		},
	}

	for _, sp := range searchPoints {
		nodesToSearch := roadindex[sp]
		minDist, closest := b.searchNodes(nodesToSearch, r)
		if minDist < minimumFoundDistance {
			minimumFoundDistance = minDist
			closestFoundNode = closest
		}
	}

	if minimumFoundDistance < BuildingThresholdDistance {
		b.ClosestRoadNode = closestFoundNode
	}
}

func (b *Building) searchNodes(nodeids []NodeId, r RoadMap) (minDistance float64, closest NodeId) {
	minDistance = 10.0

	var distance float64
	for _, nid := range nodeids {
		v := r[nid]
		a := b.AverageLocation.Longitude - v.Node.Longitude
		b := b.AverageLocation.Latitude - v.Node.Latitude
		distance = math.Hypot(a, b)
		if distance < minDistance {
			minDistance = distance
			closest = v.Node.Id
		}
	}

	return minDistance, closest
}
