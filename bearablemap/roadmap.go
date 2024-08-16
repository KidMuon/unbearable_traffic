package bearablemap

type NodeId string

type Node struct {
	Id        NodeId
	Longitude float64
	Latitude  float64
}

// Which node is the edge going to
type Edge struct {
	Id   NodeId
	Cost float64
	/* Speed and Length */
}

type RoadNode struct {
	Node  Node
	Ways  []WayId
	Edges []Edge
}

func (rn *RoadNode) AddEdge(e Edge) {
	rn.Edges = append(rn.Edges, e)
}

func (rn *RoadNode) RemoveEdge(e Edge) {
	var remaining_edges []Edge
	for _, v := range rn.Edges {
		if v == e {
			continue
		}
		remaining_edges = append(remaining_edges, v)
	}
	rn.Edges = remaining_edges
}

func (rn *RoadNode) RemoveEdgeById(id NodeId) {
	var remaining_edges []Edge
	for _, v := range rn.Edges {
		if v.Id == id {
			continue
		}
		remaining_edges = append(remaining_edges, v)
	}
	rn.Edges = remaining_edges
}

type RoadMap map[NodeId]RoadNode
type RoadIndex map[IndexPoint][]NodeId

type IndexPoint struct {
	IndexLongitude int
	IndexLatitude  int
}

func (rm RoadMap) AddRoadNode(rn RoadNode) {
	rm[rn.Node.Id] = rn
}

func (rm RoadMap) Simplify() {
	for k, v := range rm {
		if len(v.Edges) == 2 {
			rm.simplify_RemoveRoadNode(k)
		}
	}
}

func (rm RoadMap) simplify_RemoveRoadNode(id NodeId) {
	edge_a := rm[id].Edges[0]
	edge_b := rm[id].Edges[1]

	new_edge_a := Edge{
		Id:   edge_b.Id,
		Cost: edge_a.Cost + edge_b.Cost,
	}
	id_a := edge_a.Id

	new_edge_b := Edge{
		Id:   edge_a.Id,
		Cost: edge_a.Cost + edge_b.Cost,
	}
	id_b := edge_b.Id

	var change_node RoadNode
	change_node = rm[id_a]
	change_node.AddEdge(new_edge_a)
	change_node.RemoveEdgeById(id)
	rm[id_a] = change_node

	change_node = rm[id_b]
	change_node.AddEdge(new_edge_b)
	change_node.RemoveEdgeById(id)
	rm[id_b] = change_node

	delete(rm, id)
}

func (rm RoadMap) GetIndex() RoadIndex {
	index := make(RoadIndex)
	for _, rn := range rm {
		lon := int((rn.Node.Longitude - 0.005) * 100)
		lat := int((rn.Node.Latitude - 0.005) * 100)
		ip := IndexPoint{
			IndexLongitude: lon,
			IndexLatitude:  lat,
		}
		if _, ok := index[ip]; ok {
			index[ip] = append(index[ip], rn.Node.Id)
		} else {
			index[ip] = []NodeId{rn.Node.Id}
		}
	}
	return index
}

type WayMap map[WayId][]SpatialNode

type WayId string

type SpatialNode struct {
	Id          NodeId
	OrderNumber int
	SpeedLimit  int
	Longitude   float64
	Latitude    float64
}

func EliminateDisconnectedNodes(startingRoadMap RoadMap) RoadMap {
	var subsetRoadMap RoadMap
	discoveredNodes := make(map[NodeId]struct{})
	listOfSubsetRoadMap := []RoadMap{}
	var startKey NodeId
	for len(discoveredNodes) < len(startingRoadMap) {
		for k := range startingRoadMap {
			if _, ok := discoveredNodes[k]; !ok {
				startKey = k
				break
			}
		}

		subsetRoadMap = RoadMap{}
		subsetRoadMap = findConnectedNodes(startKey, subsetRoadMap, startingRoadMap)

		for k := range subsetRoadMap {
			discoveredNodes[k] = struct{}{}
		}

		listOfSubsetRoadMap = append(listOfSubsetRoadMap, subsetRoadMap)
	}

	max_subset_length := 0
	max_index := 0
	for index, subset := range listOfSubsetRoadMap {
		if len(subset) > max_subset_length {
			max_index = index
			max_subset_length = len(subset)
		}
	}

	return listOfSubsetRoadMap[max_index]
}

func findConnectedNodes(startingNode NodeId, subsetRoadMap, totalRoadMap RoadMap) RoadMap {
	subsetRoadMap[startingNode] = totalRoadMap[startingNode]
	for _, edge := range subsetRoadMap[startingNode].Edges {
		if _, ok := subsetRoadMap[edge.Id]; !ok {
			subsetRoadMap = findConnectedNodes(edge.Id, subsetRoadMap, totalRoadMap)
		}
	}
	return subsetRoadMap
}
