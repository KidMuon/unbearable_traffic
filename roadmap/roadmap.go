package roadmap

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
