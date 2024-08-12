package convert

import (
	"math"
	"strconv"

	"github.com/KidMuon/unbearable_traffic/buildingmap"
	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/roadmap"
)

func CreateRoadMap(waymap roadmap.WayMap) roadmap.RoadMap {
	roadm := make(roadmap.RoadMap)
	//First populate the roadnodes in the roadmap
	for _, sliceOfSpatialNodes := range waymap {
		for i := 0; i < len(sliceOfSpatialNodes); i++ {
			sn := sliceOfSpatialNodes[i]
			if _, ok := roadm[sn.Id]; !ok {
				roadm[sn.Id] = roadmap.RoadNode{
					Node: roadmap.Node{
						Id:        sn.Id,
						Longitude: sn.Longitude,
						Latitude:  sn.Latitude,
					},
				}
			}

			if i == 0 {
				continue
			}
			sn_prev := sliceOfSpatialNodes[i-1]
			edge_cost := distance(sn.Longitude, sn.Latitude, sn_prev.Longitude, sn_prev.Latitude)
			edge := roadmap.Edge{
				Id:   sn_prev.Id,
				Cost: edge_cost,
			}
			edge_prev := roadmap.Edge{
				Id:   sn.Id,
				Cost: edge_cost,
			}

			edge_node := roadm[sn.Id]
			edge_node.AddEdge(edge)
			roadm[sn.Id] = edge_node

			edge_prev_node := roadm[sn_prev.Id]
			edge_prev_node.AddEdge(edge_prev)
			roadm[sn_prev.Id] = edge_prev_node
		}
	}

	return roadm
}

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

func distance(lon1, lat1, lon2, lat2 float32) float32 {
	a := lon2 - lon1
	b := lat2 - lat1
	return float32(math.Hypot(float64(a), float64(b)))
}

func CreateBuildingMap(overpass_buildings overpass.OverpassBuildingResponse, road_map roadmap.RoadMap) buildingmap.BuildingMap {
	bmap := make(buildingmap.BuildingMap)

	buildings := []buildingmap.Building{}
	for _, oBuilding := range overpass_buildings.Buildings {
		bid := buildingmap.BuildingID(oBuilding.Id)

		bn := []roadmap.Node{}
		for _, bnode := range oBuilding.BuildingNodes {
			for _, node := range overpass_buildings.Nodes {
				if node.Id == bnode.Reference_id {
					longitude, _ := strconv.ParseFloat(node.Lon, 32)
					latitude, _ := strconv.ParseFloat(node.Lat, 32)

					bn = append(bn, roadmap.Node{
						Id:        roadmap.NodeId(node.Id),
						Longitude: float32(longitude),
						Latitude:  float32(latitude),
					})
					break
				}
			}
		}

		bt := buildingmap.BuildingType("")
		for _, tag := range oBuilding.BuildingTags {
			if tag.Key != "building" {
				continue
			}
			switch tag.Value {
			case "commercial", "retail":
				bt = buildingmap.CommercialBuilding
			case "leisure":
				bt = buildingmap.LeisureBuilding
			case "school":
				bt = buildingmap.SchoolBuilding
			default:
				bt = buildingmap.ResidentialBuilding
			}
		}

		b := buildingmap.Building{
			Id:              bid,
			Nodes:           bn,
			BuildingType:    bt,
			ClosestRoadNode: roadmap.NodeId(""),
		}

		b.AssignAverageLocation()
		b.AssignClosestRoadNode(road_map)

		buildings = append(buildings, b)
	}

	for _, bldg := range buildings {
		if bldg.ClosestRoadNode == roadmap.NodeId("") {
			continue
		}
		bmap[bldg.Id] = bldg
	}

	return bmap
}
