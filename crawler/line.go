package crawler

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* ========================================================================================= */
type OperationType int
type OperationTypes []OperationType

//TODO - have to map Transportation Type and Line number to OperationTypeIDMap
type OperationID string
type OperationTypeIDMap map[OperationType]OperationID

const (
	OPERATION_NORMAL OperationType = iota
	OPERATION_PRE_HOLIDAY
	OPERATION_HOLIDAY
	OPERATION_UNKNOWN
)
const (
	operation_normal_identifier                 = "делник"
	operation_pre_holiday_identifier            = "предпразник"
	operation_holiday_identifier                = "празник"
	operation_pre_holiday_or_holiday_identifier = "предпразник / празник"
)

var operationsIdentifiers = map[string]OperationType{operation_normal_identifier: OPERATION_NORMAL,
	operation_pre_holiday_or_holiday_identifier: OPERATION_HOLIDAY,
	operation_holiday_identifier:                OPERATION_HOLIDAY,
	operation_pre_holiday_identifier:            OPERATION_PRE_HOLIDAY}

var operationStrings = [...]string{OPERATION_NORMAL: "Weekday",
	OPERATION_PRE_HOLIDAY: "Pre-Holiday",
	OPERATION_HOLIDAY:     "Holiday"}

func (o OperationType) String() string {
	return operationStrings[o]
}

func getOperationsMap(doc *goquery.Document) OperationTypeIDMap {
	m := make(OperationTypeIDMap)
	doc.Find(OPERATION_ENTRY_SELECTOR).
		Each(func(i int, operationEntry *goquery.Selection) {
			typeString := strings.TrimSpace(operationEntry.Text())
			operationType := getOperationType(typeString)
			id, ok := operationEntry.Attr("id")
			if !ok {
				html, _ := operationEntry.Html()
				panic(fmt.Errorf("Operation entry MUST have 'id' element, HTML: %v", html))
			}
			//This id should be in this format at all time: schedule_{operationID}_button
			parts := strings.Split(id, "_")
			if len(parts) != 3 {
				panic(fmt.Errorf("The given id '%v' is not in the required format : 'schedule_{operationID}_button'", id))
			}
			m[operationType] = OperationID(parts[1])
		})
	return m
}

func getOperationType(operationTypeString string) OperationType {
	for identifier, operationType := range operationsIdentifiers {
		if operationTypeString == identifier {
			return operationType
		}
	}
	return OPERATION_UNKNOWN
}

/* ========================================================================================= */

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

/* ========================================================================================= */
type Line struct {
	Type   TransportationType
	Number string
	Routes
	OperationTypeIDMap
}

type Route struct {
	Direction
	Stops
}
type Routes [2]Route

type Direction struct {
	Name string
	ID   string
}
type Directions [2]Direction

type Stop struct {
	Name        string
	CapitalName string
	Sign        string
	ID          string
	URL         string
}
type Stops []Stop

/* ========================================================================================= */

const (
	LINK_SELECTOR                 = "a"
	LIST_SELECTOR                 = "li"
	OPERATION_ENTRY_SELECTOR      = ".schedule_active_list_tabs li a"
	DIRECTION_TAB_SELECTOR        = ".schedule_view_direction_tabs"
	FULL_ROUTE_SELECTOR           = ".schedule_view_route_directions"
	STOP_SIGNS_SELECTOR           = ".schedule_direction_signs"
	STOP_SIGN_SELECTOR            = ".stop_link"
	STOP_CAPITAL_NAME_SELECTOR    = ".stop_change"
	SCHEDULE_LINKS_TIMES_SELECTOR = ".schedule_times tbody a"
)

func getRoutes(doc *goquery.Document) Routes {
	//We only select one of both because the others results from their selection contains the same results
	directionSelection := doc.Find(DIRECTION_TAB_SELECTOR).First()
	fullRouteSelection := doc.Find(FULL_ROUTE_SELECTOR).First()
	directions := getDirections(directionSelection)
	oneDirectionStops, anotherDirectionStops := getRoutesStops(fullRouteSelection)
	return Routes{
		{directions[0], oneDirectionStops},
		{directions[1], anotherDirectionStops}}
}

func getDirections(directionsSelection *goquery.Selection) Directions {
	directions := Directions{}
	//There MUST be two entries and on each one we get the ID from the link
	directionsSelection.
		Find(LINK_SELECTOR).
		Each(func(i int, directionLink *goquery.Selection) {
			name := strings.TrimSpace(directionLink.Text())
			if name == "" {
				html, _ := directionLink.Html()
				panic(fmt.Errorf("Directions MUST have a name, HTML: %v", html))
			}
			url, ok := directionLink.Attr("href")
			if !ok {
				html, _ := directionLink.Html()
				panic(fmt.Errorf("Directions link MUST have a href attribute, HTML: %v", html))
			}
			splits := strings.Split(url, "/")
			id := splits[len(splits)-1]
			directions[i] = Direction{name, id}
		})
	return directions
}

func getRoutesStops(routeSelection *goquery.Selection) (Stops, Stops) {
	routesStops := make(Stops, 0)

	//We take only the first selection because they are the same regarding stops information
	routeSelection.
		Find(STOP_SIGNS_SELECTOR).
		Find(LIST_SELECTOR).
		Each(func(i int, stopSelection *goquery.Selection) {
			name, _ := getTextAndHref(stopSelection.Find(STOP_CAPITAL_NAME_SELECTOR))
			sign, url := getTextAndHref(stopSelection.Find(STOP_SIGN_SELECTOR))
			url, id := getStopURLAndID(url)
			routesStops = append(routesStops, Stop{"", name, sign, id, url})
		})
	numberOfRouteStops := len(routesStops) / 2
	return routesStops[:numberOfRouteStops], routesStops[numberOfRouteStops:]
}

func getTextAndHref(s *goquery.Selection) (string, string) {
	text := s.Text()
	href, ok := s.Attr("href")
	if !ok {
		html, _ := s.Html()
		log.Printf("The given attribute should have href attribute, HTML: %v", html)
	}
	return text, href
}

func getStopURLAndID(url string) (string, string) {
	//URL should be in this form : 'stop/{stop_id}/{stop_latin_name}#{integer_sign_name}'
	firstSlash := strings.Index(url, "/")
	lastSlash := strings.LastIndex(url, "/")
	url = url[:lastSlash]
	id := url[firstSlash+1:]
	return url, id
}

/* ========================================================================================= */

func getNormalTimes(doc *goquery.Document) []string {
	times := make([]string, 0)
	doc.Find(SCHEDULE_LINKS_TIMES_SELECTOR).
		Each(func(i int, s *goquery.Selection) {
			times = append(times, strings.TrimSpace(s.Text()))
		})
	return times
}

//func getNormalTimesOfTimes(doc *goquery.Document) [][]string {
//	timesOfTimes := make([][]string, 0)
//	doc.Find(SCHEDULE_TIMES_SELECTOR).Each(func(i1 int, st *goquery.Selection) {
//		times := make([]string, 0)
//		st.Find(SCHEDULE_LINKS_TIMES_SELECTOR).
//			Each(func(i int, s *goquery.Selection) {
//				times = append(times, strings.TrimSpace(s.Text()))
//			})
//		timesOfTimes = append(timesOfTimes, times)
//	})
//	return timesOfTimes
//}

//func advancedTimes(doc *goquery.Document) [][][]string {
//	timesOfTimes := make([][][]string, 0)
//	doc.Find(SCHEDULE_TIMES_SELECTOR).Each(func(i1 int, st *goquery.Selection) {
//		times := make([][]string, 0)
//		st.Find(SCHEDULE_LINKS_TIMES_SELECTOR).Each(func(i int, s *goquery.Selection) {
//			click, _ := s.Attr("onclick")
//			//fmt.Println(click)
//			i2 := strings.LastIndex(click, "'")
//			i1 := strings.LastIndex(click[:i2], "'")
//			reduced := click[i1+1 : i2]
//			splits := strings.Split(reduced, ",")
//			times = append(times, splits)
//		})
//		timesOfTimes = append(timesOfTimes, times)
//	})
//	return timesOfTimes
//}
//func intToTime(c string) string {
//	i, _ := strconv.Atoi(c)
//	return fmt.Sprintf("%v:%02d", i/60, i%60)
//}
//
//func printTimes(times [][]string) {
//	l := len(times[0])
//	for i := 1; i < l; i++ {
//		for _, row := range times {
//
//			time := row[i]
//
//			if time != "" {
//				fmt.Printf("%v\t", intToTime(time))
//			} else {
//				fmt.Print("*****\t")
//			}
//		}
//		fmt.Print("\n")
//	}
//}

/* ========================================================================================= */
func buildLinkSuffixes(baseURL string, something Line) []string {
	links := make([]string, 0)
	fmt.Println(something.Routes[0].Stops)
	fmt.Println(something.Routes[0].Direction.Name)
	for _, operationID := range something.OperationTypeIDMap {
		for _, route := range something.Routes {
			for _, stop := range route.Stops {
				stopID, _ := strconv.Atoi(stop.Sign)
				links = append(links,
					fmt.Sprintf("%v/%v/%v/%v", baseURL, operationID, route.Direction.ID, stopID))
			}
		}
	}
	return links
}

//func makeLineSomething(doc *goquery.Document) Line {
//	return Line{getRoutes(doc), getOperationsMap(doc)}
//}
//
//func createLinksForCrawling(baseURL string, doc *goquery.Document) []string {
//	lineSomething := makeLineSomething(doc)
//	links := buildLinkSuffixes(baseURL, lineSomething)
//	return links
//}

/* ========================================================================================= */

func CrawlLine(line LineNameAndURL, r io.Reader) {
	//doc, err := goquery.NewDocumentFromReader(r)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//links := buildLinkSuffixes(doc)
	//fmt.Println(links)
	URL := "http://schedules.sofiatraffic.bg/server/html/schedule_load"
	fmt.Println(URL)
	//links := createLinksForCrawling(URL, doc)
	//d, _ := goquery.NewDocument(links[0])
	//fmt.Println(d.Text())
	//for _, link := range links {
	//	d, _ := goquery.NewDocument(link)
	//	times := getNormalTimes(d)
	//	fmt.Printf("%v -> %v\n", link, times)
	//}

	//times := getNormalTimes(d)
	//fmt.Println(times)

	//fmt.Println(links)
	//fmt.Println(len(links))
}
