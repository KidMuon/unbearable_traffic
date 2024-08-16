package convert

import (
	"math"
	"strconv"

	"github.com/KidMuon/unbearable_traffic/bearablemap"
	"github.com/KidMuon/unbearable_traffic/overpass"
)

func CreateRoadMap(waymap bearablemap.WayMap) bearablemap.RoadMap {
	roadm := make(bearablemap.RoadMap)
	for wid, sliceOfSpatialNodes := range waymap {
		for i := 0; i < len(sliceOfSpatialNodes); i++ {
			sn := sliceOfSpatialNodes[i]
			if _, ok := roadm[sn.Id]; !ok {
				roadm[sn.Id] = bearablemap.RoadNode{
					Node: bearablemap.Node{
						Id:        sn.Id,
						Longitude: sn.Longitude,
						Latitude:  sn.Latitude,
					},
					Ways: []bearablemap.WayId{wid},
				}
			} else {
				roadm[sn.Id] = bearablemap.RoadNode{
					Node:  roadm[sn.Id].Node,
					Ways:  append(roadm[sn.Id].Ways, wid),
					Edges: roadm[sn.Id].Edges,
				}
			}

			if i == 0 {
				continue
			}
			sn_prev := sliceOfSpatialNodes[i-1]
			edge_cost := distance(sn.Longitude, sn.Latitude, sn_prev.Longitude, sn_prev.Latitude) * 4 * 69
			if sn.SpeedLimit == 0 && sn_prev.SpeedLimit == 0 {
				edge_cost /= 25.0 //Very rough conversion to get edge_cost in units of simulation time
			} else {
				edge_cost /= float64(math.Max(float64(sn.SpeedLimit), float64(sn_prev.SpeedLimit)))
			}

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
	var speedlimit int64

	type nodeLonLat struct {
		lon float64
		lat float64
	}
	nodeSpatial := make(map[string]nodeLonLat)
	for _, node := range overpass_ways.Nodes {
		loop_lon, _ := strconv.ParseFloat(node.Lon, 64)
		loop_lat, _ := strconv.ParseFloat(node.Lat, 64)
		nodeSpatial[node.Id] = nodeLonLat{
			lon: loop_lon,
			lat: loop_lat,
		}
	}

	for _, way := range overpass_ways.Streets {
		sn := []bearablemap.SpatialNode{}
		orderID := 1
		for _, wayNode := range way.StreetNodes {
			longitude = nodeSpatial[wayNode.Reference_id].lon
			latitude = nodeSpatial[wayNode.Reference_id].lat

			for _, tag := range way.StreetTags {
				if tag.Key == "maxspeed" {
					speedlimit, _ = strconv.ParseInt(tag.Value, 10, 64)
				}
			}

			sn = append(sn, bearablemap.SpatialNode{
				Id:          bearablemap.NodeId(wayNode.Reference_id),
				OrderNumber: orderID,
				SpeedLimit:  int(speedlimit),
				Longitude:   longitude,
				Latitude:    latitude,
			})
			orderID++
		}

		waymap[bearablemap.WayId(way.Id)] = sn
	}
	return waymap
}

func distance(lon1, lat1, lon2, lat2 float64) float64 {
	a := lon2 - lon1
	b := lat2 - lat1
	return float64(math.Hypot(float64(a), float64(b)))
}

func CreateBuildingMap(overpass_buildings overpass.OverpassBuildingResponse, road_map bearablemap.RoadMap) bearablemap.BuildingMap {
	bmap := make(bearablemap.BuildingMap)
	road_index := road_map.GetIndex()

	var longitude float64
	var latitude float64
	type nodeLonLat struct {
		lon float64
		lat float64
	}
	nodeSpatial := make(map[string]nodeLonLat)
	for _, node := range overpass_buildings.Nodes {
		loop_lon, _ := strconv.ParseFloat(node.Lon, 64)
		loop_lat, _ := strconv.ParseFloat(node.Lat, 64)
		nodeSpatial[node.Id] = nodeLonLat{
			lon: loop_lon,
			lat: loop_lat,
		}
	}

	buildings := []bearablemap.Building{}
	for _, oBuilding := range overpass_buildings.Buildings {
		bid := bearablemap.BuildingID(oBuilding.Id)

		bn := []bearablemap.Node{}
		for _, bnode := range oBuilding.BuildingNodes {
			longitude = nodeSpatial[bnode.Reference_id].lon
			latitude = nodeSpatial[bnode.Reference_id].lat

			bn = append(bn, bearablemap.Node{
				Id:        bearablemap.NodeId(bnode.Reference_id),
				Longitude: float64(longitude),
				Latitude:  float64(latitude),
			})

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
		b.AssignClosestRoadNode(road_map, road_index)

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
