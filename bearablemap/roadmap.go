package bearablemap

import "fmt"

type NodeId string

type Node struct {
	Id        NodeId
	Longitude float32
	Latitude  float32
}

// Which node is the edge going to
type Edge struct {
	Id   NodeId
	Cost float32
	/* Speed and Length */
}

type RoadNode struct {
	Node  Node
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

type WayMap map[WayId][]SpatialNode

type WayId string

type SpatialNode struct {
	Id          NodeId
	OrderNumber int
	Longitude   float32
	Latitude    float32
}

//I need a way to get the roads that are in the main network and eliminate those in the smaller networks
/*
Count the number of total nodes.
Get the first nodeID

While I haven't visited every node keep searching
A function that takes the nodeID that you pass to it and a small roadmap and the total roadmap
Appends the node at the nodeID to the small roadmap
Then loops over the nodes in the edges and if it finds one that isn't in the map call itself passing the nodeID and the small roadmap
	setting that to the smallroadmap it has.
If it finishes looping over the edges without finding one missing return the small roadmap.

In the main function save that smaller roadmap to a slice of roadmaps.

When the loop is finished choose the roadmap with the largest number of nodes connected within it and return that.
*/

func EliminateDisconnectedNodes(startingRoadMap RoadMap) RoadMap {
	var subsetRoadMap RoadMap
	discoveredNodes := make(map[NodeId]struct{})
	listOfSubsetRoadMap := []RoadMap{}
	safetyCounter := 0
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

		if safetyCounter > 12 {
			fmt.Println("More than 12 splits. Exiting...")
			break
		}
		safetyCounter++
	}

	max := 0
	max_index := 0
	for index, subset := range listOfSubsetRoadMap {
		//fmt.Println(len(subset))
		if len(subset) > max {
			max_index = index
			max = len(subset)
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
