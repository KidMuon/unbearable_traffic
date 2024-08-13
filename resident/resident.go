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
