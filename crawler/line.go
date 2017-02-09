package crawler

import (
	"fmt"
	"io"
	"log"
	"strings"

	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

const (
	operation_normal_identifier                 = "делник"
	operation_pre_holiday_identifier            = "предпразник"
	operation_holiday_identifier                = "празник"
	operation_pre_holiday_or_holiday_identifier = "предпразник / празник"
	operation_pre_holiday_or_holiday_prefix     = "предпразник, празник"
)

type OperationType int

const (
	OPERATION_NORMAL OperationType = iota
	OPERATION_PRE_HOLIDAY
	OPERATION_HOLIDAY
)

var operationStrings = [...]string{OPERATION_NORMAL: "Weekday",
	OPERATION_PRE_HOLIDAY: "Pre-Holiday",
	OPERATION_HOLIDAY:     "Holiday"}

func getLineOperationType(doc *goquery.Document) OperationType {
	operationTypeRaw := doc.Find(".schedule_active_list_active_tab").Text()
	switch strings.TrimSpace(operationTypeRaw) {
	case operation_normal_identifier:
		return OPERATION_NORMAL
	case operation_pre_holiday_or_holiday_identifier:
		return OPERATION_HOLIDAY
	case operation_holiday_identifier:
		return OPERATION_HOLIDAY
	case operation_pre_holiday_identifier:
		return OPERATION_PRE_HOLIDAY
	default:
		panic(fmt.Errorf("No operation mode found %v", operationTypeRaw))
	}
}

type Direction struct {
	Name string
	ID   string
	URL  string
}
type Directions [2]Direction
type StopTuple struct {
	Name         string
	Number       string
	URL          string
	StopViewLink string
}
type StopsNames []StopTuple

//ID which seems to be unique and it's for every diferent type of operation
type SCID string
type SCIDS [3]SCID

//SCID
func findDirections(doc *goquery.Document) Directions {
	directions := Directions{}
	doc.Find(".schedule_view_direction_tabs").First().Find("a").Each(func(i int, s *goquery.Selection) {

		name := strings.TrimSpace(s.Text())
		url, ok := s.Attr("href")
		if !ok {
			//TODO-Log
			fmt.Fprintf(os.Stderr, "No URL for direction view with direction with name [%v]", name)
		}
		splits := strings.Split(url, "/")
		id := splits[len(splits)-1]
		directions[i] = Direction{name, id, url}
	})
	return directions
}

func findSCIDs(doc *goquery.Document) SCIDS {
	scids := SCIDS{}
	matcher := cascadia.MustCompile(`[id^="schedule"][id$="button"]:not([id*="direction"])`)
	doc.FindMatcher(matcher).Each(func(i int, s *goquery.Selection) {
		attribute_id, _ := s.Attr("id")
		splits := strings.Split(attribute_id, "_")
		scids[i] = SCID(splits[1])
	})
	return scids
}

//get sschedule from schedule_6484_direction_172_sign_693
func findStopNames(doc *goquery.Document) StopsNames {
	stopNames := make(StopsNames, 0)
	doc.Find(".schedule_view_route_directions").
		First().
		Find(".schedule_direction_sign_wrapper ul").
		First().
		Find("li").Each(func(i int, s *goquery.Selection) {

		nameSelection := s.Find(".stop_change")
		name := nameSelection.Text()
		stopViewLink, ok := nameSelection.Attr("href")
		if !ok {
			//TODO-Log
			fmt.Fprintf(os.Stderr, "No URL for stop view for stop with name [%v]", name)
		}

		numberSelection := s.Find(".stop_link")
		number := numberSelection.Text()
		url, ok := numberSelection.Attr("href")
		if !ok {
			//TODO-Log
			fmt.Fprintf(os.Stderr, "No URL for stop with name [%v] and number [%v]", name, number)
		}

		stopNames = append(stopNames, StopTuple{name, number, url, stopViewLink})
	})

	return stopNames
}

func getSelectionForOperationType(doc *goquery.Document, operationType OperationType) {
	doc.Find(".schedule_active_list_content").Each(func(i int, operationTypeSelection *goquery.Selection) {
		str := strings.TrimSpace(operationTypeSelection.ChildrenFiltered("h3").Text())
		var operationType OperationType
		switch true {
		case strings.HasPrefix(str, string(operation_normal_identifier)):
			operationType = OPERATION_NORMAL
		case strings.HasPrefix(str, string(operation_pre_holiday_or_holiday_prefix)):
			operationType = OPERATION_HOLIDAY
		case strings.HasPrefix(str, string(operation_holiday_identifier)):
			operationType = OPERATION_HOLIDAY
		case strings.HasPrefix(str, string(operation_pre_holiday_identifier)):
			operationType = OPERATION_PRE_HOLIDAY
		}

		fmt.Printf("Operation type: [%v]\n", operationStrings[operationType])

	})

}

func CrawlLine(line LineNameAndURL, r io.Reader) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	getSelectionForOperationType(doc, OPERATION_HOLIDAY)
	directions := findDirections(doc)
	fmt.Println(directions)
	stopsNames := findStopNames(doc)
	fmt.Println(stopsNames)
	scids := findSCIDs(doc)
	fmt.Println(scids)
}
