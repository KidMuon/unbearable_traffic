package bearabletime

import "math/rand"

const simulationsections int = 96 //Tick every 15 minutes

func RandomStartTime(shift int) int {
	startTime := 0
	switch shift {
	case 2:
		startTime = 13*4 + rand.Intn(3*4)
	case 3:
		startTime = 18*4 + rand.Intn(5*4)
	default:
		//assume first shift unless a valid value is given
		startTime = 7*4 + rand.Intn(2*4)
	}
	if startTime >= simulationsections {
		startTime -= simulationsections
	}
	return startTime
}

func RandomEndTime(shift int) int {
	endTime := 0
	switch shift {
	case 2:
		endTime = 21*4 + rand.Intn(2*4)
	case 3:
		endTime = 6*4 + rand.Intn(3*4)
	default:
		//assume first shift unless a valid value is given
		endTime = 17*4 + rand.Intn(5*4)
	}
	if endTime >= simulationsections {
		endTime -= simulationsections
	}
	return endTime
}

func RegularSchedule(locations, startTime, endTime int) []int {
	if endTime < startTime {
		endTime += simulationsections
	}

	schedule := []int{startTime}
	splits := locations - 1
	for i := 1; i < splits; i++ {
		nextTime := startTime + ((endTime-startTime)/splits)*i
		schedule = append(schedule, nextTime)
	}
	schedule = append(schedule, endTime)

	for i := range schedule {
		if schedule[i] >= simulationsections {
			schedule[i] -= simulationsections
		}
	}

	return schedule
}
