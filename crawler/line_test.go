package crawler

import (
	"strings"
	"testing"
)

func TestOperationIDRoutesMap_String(t *testing.T) {
	operationIDRoutesMap := OperationIDRoutesMap{
		"4000": Routes{
			Route{
				Direction{Name: "First", ID: "1"},
				Stops{Stop{Sign: "100", CapitalName: "STATION 1", ID: "30"}},
			},
			Route{
				Direction{Name: "Second", ID: "2"},
				Stops{Stop{Sign: "110", CapitalName: "STATION 2", ID: "31"}}},
		},
		"4001": Routes{
			Route{
				Direction{Name: "First", ID: "1"},
				Stops{Stop{Sign: "100", CapitalName: "STATION 1", ID: "30"}},
			}},
	}
	firstRoutesID := "(OperationID: 4000)\nFirst (1)\n[100] STATION 1 (30)\n\nSecond (2)\n[110] STATION 2 (31)\n\n"
	secondRoutesID := "(OperationID: 4001)\nFirst (1)\n[100] STATION 1 (30)\n\n"

	s := operationIDRoutesMap.String()
	if !strings.Contains(s, firstRoutesID) {
		t.Errorf(`"%#v".String() should contain %q, but it doesnt %q`, operationIDRoutesMap, firstRoutesID, s)
	}
	if !strings.Contains(s, secondRoutesID) {
		t.Errorf(`"%#v".String() should contain %q, but it doesnt %q`, operationIDRoutesMap, secondRoutesID, s)
	}
}

func TestLine_String(t *testing.T) {
	line := Line{
		Name:           "10",
		Transportation: Tram,
		OperationIDMap: OperationIDMap{
			Normal: "4000",
		},
		OperationIDRoutesMap: OperationIDRoutesMap{
			"4000": Routes{
				Route{
					Direction{Name: "First", ID: "1"},
					Stops{Stop{Sign: "100", CapitalName: "STATION 1", ID: "30"}},
				},
				Route{
					Direction{Name: "Second", ID: "2"},
					Stops{Stop{Sign: "110", CapitalName: "STATION 2", ID: "31"}}},
			}},
	}
	expected := "Tram '10'\nWeekday (4000)\n(OperationID: 4000)\nFirst (1)\n[100] STATION 1 (30)\n\nSecond (2)\n[110] STATION 2 (31)\n\n"

	s := line.String()
	if s != expected {
		t.Errorf(`"%#v".String() must be => %q, given %q`, line, expected, s)
	}
}

func TestLine_scheduleIDs(t *testing.T) {
	line := Line{
		Name:           "10",
		Transportation: Tram,
		OperationIDMap: OperationIDMap{
			Normal: "4000",
		},
		OperationIDRoutesMap: OperationIDRoutesMap{
			"4000": Routes{
				Route{
					Direction{Name: "First", ID: "1"},
					Stops{
						Stop{Sign: "100", CapitalName: "STATION 1", ID: "30"},
						Stop{Sign: "0910", CapitalName: "STATION 3", ID: "30"},
					},
				},
				Route{
					Direction{Name: "Second", ID: "2"},
					Stops{Stop{Sign: "0211", CapitalName: "STATION 2", ID: "31"}}},
			}},
	}
	expected := []ScheduleID{
		ScheduleID("4000/1/100"),
		ScheduleID("4000/1/910"),
		ScheduleID("4000/2/211"),
	}
	result := line.scheduleIDs()

	if len(result) != len(expected) {
		t.Errorf(`line.scheduleIDs() produces a different amount of IDs, expected: %q, got: %q`, expected, result)
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf(`line.scheduleIDs() expects at %d-th place => %q, got %q`, i, expected[i], result[i])
		}
	}
}
