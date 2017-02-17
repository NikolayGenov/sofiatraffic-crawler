package crawler

import (
	"fmt"
	"testing"
)

func TestStop_String(t *testing.T) {
	stop := Stop{
		Sign:             "0910",
		CapitalName:      "SOME STATION",
		ID:               "1010",
		Name:             "Some station",
		VirtualTableStop: VirtualTableStop{"1", "2", "3", "4"},
	}
	expected := "[0910] SOME STATION (1010)"
	s := fmt.Sprintf("%s", stop.String())
	if s != expected {
		t.Errorf(`"%v".String() must be => %q, given %q`, stop, expected, s)
	}
}

func TestStops_StringEmpty(t *testing.T) {
	stops := Stops{}
	s := fmt.Sprintf("%s", stops.String())
	if s != "" {
		t.Errorf(`"%v".String() must be empty string, given %q`, stops, s)
	}
}

func TestStops_String(t *testing.T) {
	stops := Stops{
		Stop{
			Sign:             "0910",
			CapitalName:      "SOME STATION",
			ID:               "1010",
			Name:             "Some station",
			VirtualTableStop: VirtualTableStop{"1", "2", "3", "4"}},
		Stop{
			Sign:             "2313",
			CapitalName:      "SOME OTHER STATION",
			ID:               "1022",
			Name:             "Some other station",
			VirtualTableStop: VirtualTableStop{"1", "2", "3", "5"},
		},
	}
	expected := "[0910] SOME STATION (1010)\n[2313] SOME OTHER STATION (1022)\n"
	s := fmt.Sprintf("%s", stops.String())
	if s != expected {
		t.Errorf(`"%v".String() must be %q, given %q`, stops, expected, s)
	}
}

func TestRoute_String(t *testing.T) {
	route :=
		Route{
			Direction{
				Name: "Liberty - j.k. Last Hope",
				ID:   "3451",
			},
			Stops{
				Stop{
					Sign:             "0910",
					CapitalName:      "SOME STATION",
					ID:               "1010",
					Name:             "Some station",
					VirtualTableStop: VirtualTableStop{"1", "2", "3", "4"}},
				Stop{
					Sign:             "2313",
					CapitalName:      "SOME OTHER STATION",
					ID:               "1022",
					Name:             "Some other station",
					VirtualTableStop: VirtualTableStop{"1", "2", "3", "5"},
				}}}
	expected := "\nLiberty - j.k. Last Hope (3451)\n[0910] SOME STATION (1010)\n[2313] SOME OTHER STATION (1022)\n"
	s := fmt.Sprintf("%s", route.String())
	if s != expected {
		t.Errorf(`"%v".String() must be %q, given %q`, route, expected, s)
	}
}

func TestRoutes_String(t *testing.T) {
	routes := Routes{
		Route{
			Direction{
				Name: "Liberty - j.k. Last Hope",
				ID:   "3451",
			},
			Stops{
				Stop{
					Sign:             "0910",
					CapitalName:      "SOME STATION",
					ID:               "1010",
					Name:             "Some station",
					VirtualTableStop: VirtualTableStop{"1", "2", "3", "4"}},
				Stop{
					Sign:             "2313",
					CapitalName:      "SOME OTHER STATION",
					ID:               "1022",
					Name:             "Some other station",
					VirtualTableStop: VirtualTableStop{"1", "2", "3", "5"},
				}}},
		Route{
			Direction{
				Name: "j.k. Last Hope - Liberty",
				ID:   "3452",
			},
			Stops{
				Stop{
					Sign:             "2314",
					CapitalName:      "SOME OTHER STATION",
					ID:               "1122",
					Name:             "Some other station",
					VirtualTableStop: VirtualTableStop{"1", "6", "1", "10"}},
				Stop{
					Sign:             "0911",
					CapitalName:      "SOME STATION",
					ID:               "1100",
					Name:             "Some station",
					VirtualTableStop: VirtualTableStop{"1", "6", "1", "9"}},
			}}}
	expected := "\nLiberty - j.k. Last Hope (3451)\n[0910] SOME STATION (1010)\n[2313] SOME OTHER STATION (1022)\n" +
		"\nj.k. Last Hope - Liberty (3452)\n[2314] SOME OTHER STATION (1122)\n[0911] SOME STATION (1100)\n"
	s := fmt.Sprintf("%s", routes.String())
	if s != expected {
		t.Errorf(`"%v".String() must be %q, given %q`, routes, expected, s)
	}
}
