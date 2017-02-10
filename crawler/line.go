package crawler

import (
	"fmt"
	"io"
	"log"
	"strings"

	"os"

	"github.com/PuerkitoBio/goquery"
)

const (
	operation_normal_identifier                 = "делник"
	operation_pre_holiday_identifier            = "предпразник"
	operation_holiday_identifier                = "празник"
	operation_pre_holiday_or_holiday_identifier = "предпразник / празник"
	operation_pre_holiday_or_holiday_prefix     = "предпразник, празник"
)

type OperationType int

type OperationTypes []OperationType

const (
	OPERATION_NORMAL OperationType = iota
	OPERATION_PRE_HOLIDAY
	OPERATION_HOLIDAY
	OPERATION_UNKNOWN
)

var operationsIdentifiersMap = map[string]OperationType{operation_normal_identifier: OPERATION_NORMAL,
	operation_pre_holiday_or_holiday_identifier: OPERATION_HOLIDAY,
	operation_holiday_identifier:                OPERATION_HOLIDAY,
	operation_pre_holiday_identifier:            OPERATION_PRE_HOLIDAY}

var operationStrings = [...]string{OPERATION_NORMAL: "Weekday",
	OPERATION_PRE_HOLIDAY: "Pre-Holiday",
	OPERATION_HOLIDAY:     "Holiday"}

func (o OperationType) String() string {
	return operationStrings[o]
}

// ===============================================
//Current date and page postion related functions
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

// =============================================
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
type OperationID string
type OperationTypeIDMap map[OperationType]OperationID

//
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

//get sschedule from schedule_6484_direction_172_sign_693
func findStopNames(doc *goquery.Document) (StopsNames, StopsNames) {
	stopNames := make(StopsNames, 0)
	doc.Find(".schedule_view_route_directions").
		First().
		Find(".schedule_direction_sign_wrapper ul").
		//First().
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

	return stopNames[:len(stopNames)/2], stopNames[len(stopNames)/2:]
}
func findOperationType(s string) OperationType {
	for k, v := range operationsIdentifiersMap {
		if s == k {
			return v
		}
	}
	return OPERATION_UNKNOWN
}

func getOperationMap(doc *goquery.Document) OperationTypeIDMap {
	operationMap := make(OperationTypeIDMap)
	doc.Find(".schedule_active_list_tabs li a").Each(func(i int, linkSelection *goquery.Selection) {
		operationTypeString := strings.TrimSpace(linkSelection.Text())
		operationType := findOperationType(operationTypeString)
		if attributeID, ok := linkSelection.Attr("id"); ok {
			attributeParts := strings.Split(attributeID, "_")
			if len(attributeParts) > 1 {
				operationID := attributeParts[1]
				operationMap[operationType] = OperationID(operationID)
			} else {
				log.Printf("id of this link is not in the format 'schedule_xxxx_button' as needed: %v", attributeID)
			}
		} else {
			log.Println("This element does not have attribute id, which was expected for normal processing")

		}

	})
	return operationMap
}
func getTimes(doc *goquery.Document) []string {
	//TODO - Only One selection
	times := make([]string, 0)
	doc.Find(".schedule_times tbody").First().Find("a").Each(func(i int, s *goquery.Selection) {
		times = append(times, strings.TrimSpace(s.Text()))
	})
	return times
}

//func advancedTimes(){
//}

func CrawlLine(line LineNameAndURL, r io.Reader) {
	//doc, err := goquery.NewDocument("http://schedules.sofiatraffic.bg/server/html/schedule_load/6672/2696/377")
	//fmt.Println(doc.Text())
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	m := getOperationMap(doc)
	fmt.Println(m)
	directions := findDirections(doc)
	fmt.Println(directions)
	stopsNames1, stopnames2 := findStopNames(doc)
	fmt.Println(stopsNames1)
	fmt.Println(stopnames2)
	times := getTimes(doc)
	fmt.Println(times)

}
