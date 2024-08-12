package convert

import (
	"math"
	"strconv"

	"github.com/KidMuon/unbearable_traffic/bearablemap"
	"github.com/KidMuon/unbearable_traffic/overpass"
)

func CreateRoadMap(waymap bearablemap.WayMap) bearablemap.RoadMap {
	roadm := make(bearablemap.RoadMap)
	//First populate the roadnodes in the roadmap
	for _, sliceOfSpatialNodes := range waymap {
		for i := 0; i < len(sliceOfSpatialNodes); i++ {
			sn := sliceOfSpatialNodes[i]
			if _, ok := roadm[sn.Id]; !ok {
				roadm[sn.Id] = bearablemap.RoadNode{
					Node: bearablemap.Node{
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
			edge := bearablemap.Edge{
				Id:   sn_prev.Id,
				Cost: edge_cost,
			}
			edge_prev := bearablemap.Edge{
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

func CreateWayMap(overpass_ways overpass.OverpassStreetResponse) bearablemap.WayMap {
	waymap := make(bearablemap.WayMap)
	var longitude float64
	var latitude float64
	for _, way := range overpass_ways.Streets {
		sn := []bearablemap.SpatialNode{}
		orderID := 1
		for _, wayNode := range way.StreetNodes {
			for _, node := range overpass_ways.Nodes {
				if node.Id == wayNode.Reference_id {
					longitude, _ = strconv.ParseFloat(node.Lon, 32)
					latitude, _ = strconv.ParseFloat(node.Lat, 32)
				}
			}
			sn = append(sn, bearablemap.SpatialNode{
				Id:          bearablemap.NodeId(wayNode.Reference_id),
				OrderNumber: orderID,
				Longitude:   float32(longitude),
				Latitude:    float32(latitude),
			})
			orderID++
		}

		waymap[bearablemap.WayId(way.Id)] = sn
	}
	return waymap
}

func distance(lon1, lat1, lon2, lat2 float32) float32 {
	a := lon2 - lon1
	b := lat2 - lat1
	return float32(math.Hypot(float64(a), float64(b)))
}

func CreateBuildingMap(overpass_buildings overpass.OverpassBuildingResponse, road_map bearablemap.RoadMap) bearablemap.BuildingMap {
	bmap := make(bearablemap.BuildingMap)

	buildings := []bearablemap.Building{}
	for _, oBuilding := range overpass_buildings.Buildings {
		bid := bearablemap.BuildingID(oBuilding.Id)

		bn := []bearablemap.Node{}
		for _, bnode := range oBuilding.BuildingNodes {
			for _, node := range overpass_buildings.Nodes {
				if node.Id == bnode.Reference_id {
					longitude, _ := strconv.ParseFloat(node.Lon, 32)
					latitude, _ := strconv.ParseFloat(node.Lat, 32)

					bn = append(bn, bearablemap.Node{
						Id:        bearablemap.NodeId(node.Id),
						Longitude: float32(longitude),
						Latitude:  float32(latitude),
					})
					break
				}
			}
		}

		bt := bearablemap.BuildingType("")
		for _, tag := range oBuilding.BuildingTags {
			if tag.Key != "building" {
				continue
			}
			switch tag.Value {
			case "commercial", "retail":
				bt = bearablemap.CommercialBuilding
			case "leisure":
				bt = bearablemap.LeisureBuilding
			case "school":
				bt = bearablemap.SchoolBuilding
			default:
				bt = bearablemap.ResidentialBuilding
			}
		}

		b := bearablemap.Building{
			Id:              bid,
			Nodes:           bn,
			BuildingType:    bt,
			ClosestRoadNode: bearablemap.NodeId(""),
		}

		b.AssignAverageLocation()
		b.AssignClosestRoadNode(road_map)

		buildings = append(buildings, b)
	}

	for _, bldg := range buildings {
		if bldg.ClosestRoadNode == bearablemap.NodeId("") {
			continue
		}
		bmap[bldg.Id] = bldg
	}

	return bmap
}
