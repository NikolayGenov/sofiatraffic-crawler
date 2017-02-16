package crawler

import "fmt"

//Transportation represents all possible types of transportation that are supported
type Transportation int

//Note that the order is not random and Tram should have Transportation = 0, Bus = 1 and Trolley = 2
// Because those integers are used by Virtual Tables site as ids for queries
// Also Subway is not supported for now
const (
	//Tram is representing all tramway lines
	Tram Transportation = iota

	//Bus is representing all urban bus lines and all suburban bus lines
	Bus

	//Trolleybus is representing all trolleybus transportation
	Trolley
)

var (
	transportationIdentifier = map[string]Transportation{
		"tramway":    Tram,
		"trolleybus": Trolley,
		"autobus":    Bus}

	transportationStrings = [...]string{
		Tram:    "Tram",
		Bus:     "Bus",
		Trolley: "Trolley"}
)

func (t Transportation) String() string {
	return transportationStrings[t]
}

//convertToTransportation it takes a identifier (could be from an URL) and parses that to Transportation type
func convertToTransportation(identifier string) (Transportation, error) {
	t, ok := transportationIdentifier[identifier]
	if !ok {
		return -1, fmt.Errorf("Unrecognized identifer for Transporation type: %v", identifier)
	}
	return t, nil
}
