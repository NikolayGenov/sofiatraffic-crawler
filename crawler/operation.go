package crawler

import "fmt"

//Operation is used to denote different Line Operation modes [Normal, Pre Holiday, Holiday].
type Operation int

//OperationID is unique ID taken from schedules.sofiatraffic.bg, which is a combination of line
// and its Operation mode. One line can have multiple OperationIDs so we use OperationIDMap.
type OperationID string

//OperationIDMap is mapping between Operation and OperationID.
// This is needed because each line has different number of Operation modes
// raging from 0 modes to 3 different modes and each one has unique OperationID.
type OperationIDMap map[Operation]OperationID

const (
	//Normal operation mode denotes regular everyday operation of the line,normally Weekdays.
	// But it also can be all week.
	Normal Operation = iota

	//PreHoliday operation mode denotes a period (usually a day) before big holiday or usually Saturday if
	// there is a different times of operation for Holiday.
	PreHoliday

	//Holiday operation mode denotes holiday mode of operation and times.
	// Usually Sunday, but can be any official holiday, also used when there is no difference in
	// time of operation with PreHoliday.
	Holiday
)

var (
	//Mapping between identifiers found on a line page and Operation mode that will be used
	// "делник" - Weekday , "предпразник" - PreHoliday , "празник" - Holiday
	operationsIdentifiers = map[string]Operation{
		"делник":                         Normal,
		"предпразник":                    PreHoliday,
		"празник":                        Holiday,
		"предпразник / празник":          Holiday,
		"делник / предпразник / празник": Normal}

	operationStrings = [...]string{
		Normal:     "Weekday",
		PreHoliday: "Pre-Holiday",
		Holiday:    "Holiday"}
)

func (o Operation) String() string {
	return operationStrings[o]
}

func (o *OperationIDMap) String() string {
	s := ""
	if len(*o) == 0 {
		return "Line is not operational"
	}
	for operation, id := range *o {
		s += fmt.Sprintf("%v (%v)\n", operation, id)
	}
	return s
}

//convertToOperation uses operationsIdentifiers to map convert UTF8 local language strings to Operation type
func convertToOperation(identifier string) (Operation, error) {
	t, ok := operationsIdentifiers[identifier]
	if !ok {
		return -1, fmt.Errorf("Unrecognized identifer for Operation type: %v", identifier)
	}
	return t, nil
}
