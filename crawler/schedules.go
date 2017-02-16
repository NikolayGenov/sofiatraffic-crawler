package crawler

import "strings"

type ScheduleID string

type ScheduleTimes []string

//type Schedules map[ScheduleID]ScheduleTimes

func convertToScheduleID(path string) ScheduleID {
	//The path is in this format /server/html/schedule_load/{OperationID}/{DirectionID}/{StopSign}
	//ScheduleID is in this format {OperationID}/{DirectionID}/{StopSign}
	parts := strings.Split(path, "/")
	return ScheduleID(strings.Join(parts[len(parts)-3:], "/"))
}
