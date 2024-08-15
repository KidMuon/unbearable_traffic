package resident

import (
	"math/rand"

	"github.com/KidMuon/unbearable_traffic/bearablemap"
	"github.com/KidMuon/unbearable_traffic/bearabletime"
)

type Person struct {
	CurrentLocationIndex int
	LocationsToVisit     []bearablemap.NodeId
	WorkingShift         int
	DayStart             int
	DayEnd               int
	DepartureTimes       []int
}

func CreatePopulation(n int, structuremap bearablemap.BuildingMap) []Person {
	population := []Person{}

	//create an array that will be used to determine shift
	//shifts are unequally distributed
	shiftRatios := []int{100, 25, 10}
	shift := []int{}
	for i, count := range shiftRatios {
		for j := 0; j < count; j++ {
			shift = append(shift, i+1)
		}
	}

	//number of locations to visit is unequally distributed
	locationRatios := []int{1000, 800, 500, 100, 50}
	location := []int{}
	for i, count := range locationRatios {
		for j := 0; j < count; j++ {
			location = append(location, i+2)
		}
	}

	for i := 0; i < n; i++ {
		resident := Person{
			CurrentLocationIndex: 0,
			LocationsToVisit:     []bearablemap.NodeId{},
		}

		tovisit := location[rand.Intn(len(location))]
		for _, idx := range rand.Perm(len(structuremap))[0:tovisit] {
			i := 0
			for _, building := range structuremap {
				if i == idx {
					nextlocationtovisit := building.ClosestRoadNode
					resident.LocationsToVisit = append(resident.LocationsToVisit, nextlocationtovisit)
					break
				} else {
					i++
				}
			}
		}

		resident.WorkingShift = shift[rand.Intn(len(shift))]
		resident.DayStart = bearabletime.RandomStartTime(resident.WorkingShift)
		resident.DayEnd = bearabletime.RandomEndTime(resident.WorkingShift)
		resident.DepartureTimes = bearabletime.RegularSchedule(len(resident.LocationsToVisit), resident.DayStart, resident.DayEnd)

		population = append(population, resident)
	}
	return population
}

func (p Person) ReadyToTravel(time int) bool {
	return p.DepartureTimes[p.CurrentLocationIndex] == time
}

func (p *Person) Travel() {
	p.CurrentLocationIndex++
	if p.CurrentLocationIndex >= len(p.LocationsToVisit) {
		p.CurrentLocationIndex = 0
	}
}

type Route struct {
	Sections []RouteSection
}

type RouteSection struct {
	StreetNumber bearablemap.WayId
	FromNode     bearablemap.NodeId
	ToNode       bearablemap.NodeId
	Time         int
}

type AstarRouteMap map[bearablemap.NodeId]AstarRow

type AstarRow struct {
	F    float64
	G    float64
	H    float64
	Last bearablemap.NodeId
}

type openSet map[bearablemap.NodeId]float64

func (p Person) CalculateRoute(roadmap bearablemap.RoadMap, startTime int) Route {
	optimalRoute := Route{}

	var startNode, endNode bearablemap.NodeId
	startNode = p.LocationsToVisit[p.CurrentLocationIndex]
	if p.CurrentLocationIndex >= len(p.LocationsToVisit) {
		endNode = p.LocationsToVisit[0]
	} else {
		endNode = p.LocationsToVisit[p.CurrentLocationIndex+1]
	}

	if startNode == endNode {
		return optimalRoute
	}

	openset := make(openSet)
	routemap := make(AstarRouteMap)
	routemap[startNode] = AstarRow{
		F:    0.0,
		G:    0.0,
		H:    0.0,
		Last: "",
	}

	optimalRoute = AStar(startNode, endNode, roadmap, routemap, openset)

	for i := 0; i < len(optimalRoute.Sections); i++ {
		optimalRoute.Sections[i].Time += startTime
	}

	return optimalRoute
}

func AStar(current, end bearablemap.NodeId, roadmap bearablemap.RoadMap, routemap AstarRouteMap, openset openSet) Route {

	var route Route

	for _, edge := range roadmap[current].Edges {
		nid := edge.Id
		if nid == routemap[current].Last {
			continue
		}

		g := routemap[current].G + float64(edge.Cost)
		h := distanceBetweenNodes(roadmap[current].Node, roadmap[nid].Node)
		f := g + h
		if _, ok := routemap[nid]; !ok || routemap[nid].F > f {
			routemap[nid] = AstarRow{
				F:    f,
				G:    g,
				H:    h,
				Last: current,
			}
			openset[nid] = f
		}
	}

	var next bearablemap.NodeId
	min_f := 1000.0
	for oid, f := range openset {
		if f < min_f {
			min_f = f
			next = oid
		}
	}

	if next != end {
		delete(openset, next)
		route = AStar(next, end, roadmap, routemap, openset)
	} else {
		toNode := end //before the loop set toNode to Next
		for {
			fromNode := routemap[toNode].Last
			if fromNode == bearablemap.NodeId("") {
				break
			}

			wayid := bearablemap.WayId("")
			for _, wi := range roadmap[fromNode].Ways {
				for _, wj := range roadmap[toNode].Ways {
					if wi == wj {
						wayid = wi
					}
				}
			}

			section := RouteSection{
				StreetNumber: wayid,
				FromNode:     fromNode,
				ToNode:       toNode,
				Time:         int(routemap[toNode].F),
			}

			route.Sections = append(route.Sections, section)
			toNode = fromNode
		}
	}

	return route
}

func distanceBetweenNodes(m, n bearablemap.Node) float64 {
	return 1.0
}
