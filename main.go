package main

import (
	"fmt"
	"log"

	"github.com/KidMuon/unbearable_traffic/overpass"
	"github.com/KidMuon/unbearable_traffic/roadmap"
)

func main() {
	overpassData, err := overpass.ImportOverpassData_Standard()
	if err != nil {
		log.Fatal(err)
	}
	overpass.SummarizeOverpassData(overpassData)

	rn_A := roadmap.RoadNode{
		Node: roadmap.Node{
			Id: "A",
		},
		Edges: []roadmap.Edge{
			{
				Id:   "B",
				Cost: 10,
			},
		},
	}

	rn_B := roadmap.RoadNode{
		Node: roadmap.Node{
			Id: "B",
		},
		Edges: []roadmap.Edge{
			{
				Id:   "A",
				Cost: 10,
			},
			{
				Id:   "C",
				Cost: 10,
			},
		},
	}

	rn_C := roadmap.RoadNode{
		Node: roadmap.Node{
			Id: "C",
		},
		Edges: []roadmap.Edge{
			{
				Id:   "B",
				Cost: 10,
			},
			{
				Id:   "D",
				Cost: 10,
			},
			{
				Id:   "E",
				Cost: 10,
			},
		},
	}

	rn_D := roadmap.RoadNode{
		Node: roadmap.Node{
			Id: "D",
		},
		Edges: []roadmap.Edge{
			{
				Id:   "C",
				Cost: 10,
			},
		},
	}

	rn_E := roadmap.RoadNode{
		Node: roadmap.Node{
			Id: "E",
		},
		Edges: []roadmap.Edge{
			{
				Id:   "C",
				Cost: 10,
			},
		},
	}

	rm := make(roadmap.RoadMap)
	rm["A"] = rn_A
	rm["B"] = rn_B
	rm["C"] = rn_C
	rm["D"] = rn_D
	rm["E"] = rn_E

	fmt.Println(rm)
	rm.Simplify()
	fmt.Println(rm)
}
