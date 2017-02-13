package crawler

import "fmt"

type Transportation int

const (
	Tram Transportation = iota
	Bus
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

func convertToTransportation(identifier string) (Transportation, error) {
	t, ok := transportationIdentifier[identifier]
	if !ok {
		return -1, fmt.Errorf("Unrecognized identifer for Transporation type: %v", identifier)
	}
	return t, nil
}
