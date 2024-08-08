package convert

import (
	"strconv"

	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/roadmap"
)

func CreateWayMap(overpass_ways overpass.OverpassStreetResponse) roadmap.WayMap {
	waymap := make(roadmap.WayMap)
	var longitude float64
	var latitude float64
	for _, way := range overpass_ways.Streets {
		sn := []roadmap.SpatialNode{}
		orderID := 1
		for _, wayNode := range way.StreetNodes {
			for _, node := range overpass_ways.Nodes {
				if node.Id == wayNode.Reference_id {
					longitude, _ = strconv.ParseFloat(node.Lon, 32)
					latitude, _ = strconv.ParseFloat(node.Lat, 32)
				}
			}
			sn = append(sn, roadmap.SpatialNode{
				Id:          roadmap.NodeId(wayNode.Reference_id),
				OrderNumber: orderID,
				Longitude:   float32(longitude),
				Latitude:    float32(latitude),
			})
			orderID++
		}
		waymap[roadmap.WayId(way.Id)] = sn
	}
	return waymap
}

func CreateRoadMap(waymap roadmap.WayMap) roadmap.RoadMap {
	roadmap := make(roadmap.RoadMap)
	//First populate the roadnodes in the roadmap
	for _, sliceOfSpatialNodes := range waymap {
		for _, sn := range sliceOfSpatialNodes {
			rn := roadmap.RoadNode{
				Node: roadmap.Node{
					Id:        sn.Id,
					Longitude: sn.Longitude,
					Latitude:  sn.Latitude,
				},
			}
			roadmap[sn.Id] = rn
		}
	}
	//Second go through creating edges
	return roadmap
}
