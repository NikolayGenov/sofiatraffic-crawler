package crawler

import (
	"fmt"
	"strings"
	"testing"
)

var operationTests = []struct {
	in  Operation
	out string
}{
	{Normal, "Weekday"},
	{PreHoliday, "Pre-Holiday"},
	{Holiday, "Holiday"},
}

func TestOperation_String(t *testing.T) {
	for _, tt := range operationTests {
		s := fmt.Sprintf("%s", tt.in.String())
		if s != tt.out {
			t.Errorf(`"%v".String() must be => %q, given %q`, tt.in, tt.out, s)
		}
	}
}

var operationIDMapTests = []struct {
	in  OperationIDMap
	out []string
}{
	{OperationIDMap{},
		[]string{"Line is not operational"}},
	{OperationIDMap{Normal: "1011"},
		[]string{"Weekday (1011)\n"}},
	{OperationIDMap{Holiday: "11"},
		[]string{"Holiday (11)\n"}},
	{OperationIDMap{PreHoliday: "12"},
		[]string{"Pre-Holiday (12)"}},
	{OperationIDMap{Holiday: "13", PreHoliday: "12"},
		[]string{"Pre-Holiday (12)\n", "Holiday (13)\n"}},
	{OperationIDMap{Normal: "11", Holiday: "13", PreHoliday: "12"},
		[]string{"Pre-Holiday (12)\n", "Weekday (11)\n", "Holiday (13)\n"}},
}

func TestOperationIDMap_String(t *testing.T) {
	for _, tt := range operationIDMapTests {
		s := fmt.Sprintf("%s", tt.in.String())
		for _, str := range tt.out {
			if !strings.Contains(s, str) {
				t.Errorf(`"%v".String() must constain => %q,  given %q`, tt.in, str, s)
			}
		}
	}
}

var operationConversionTests = []struct {
	in  string
	out Operation
	err string
}{
	{"делник", Normal, ""},
	{"предпразник", PreHoliday, ""},
	{"празник", Holiday, ""},
	{"предпразник / празник", Holiday, ""},
	{"делник / предпразник / празник", Normal, ""},
	{"UNKNOWN", -1, "Unrecognized identifer for Operation type: UNKNOWN"},
}

func TestConvertToOperation(t *testing.T) {
	for _, tt := range operationConversionTests {
		op, err := convertToOperation(tt.in)

		if err != nil && err.Error() != tt.err {
			t.Errorf(`Converting from '%s' should produce error =>  %q, but it returned %q`, tt.in, tt.err, err)
		} else if op != tt.out && err == nil {
			t.Errorf(`Converting from '%s' must return => (%q, %q), but it returned (%q, %q)`, tt.in, tt.out, tt.err, op, err)
		}
	}

}
