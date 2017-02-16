package crawler

import "strings"

//ScheduleID is string in this format : {OperationID}/{DirectionID}/(without leading zeros){Stop.Sign}
// It can be used to query information about the stop schedule from an internal server
type ScheduleID string

//ScheduleTimes is a list of times in x:XX time string format e.g [5:13 6:49 10:23 23:01]
type ScheduleTimes []string

func convertToScheduleID(path string) ScheduleID {
	// The path is in this format /server/html/schedule_load/{OperationID}/{DirectionID}/{StopSign}
	// ScheduleID is in this format {OperationID}/{DirectionID}/{StopSign}
	parts := strings.Split(path, "/")
	return ScheduleID(strings.Join(parts[len(parts)-3:], "/"))
}
