package main

import (
	"fmt"

	"../crawler"
)

func main() {

	st, err := crawler.NewSofiaTrafficCrawler(":6379")
	if err != nil {
		panic(err)
	}

	st.CrawlLines()
	//for _, line := range st.Lines {
	//	fmt.Printf("Line [%v %v]\nOperations:  %v\n", line.Transportation, line.Name, &line.OperationIDMap)
	//	for operationID, routes := range line.OperationIDRoutesMap {
	//		fmt.Printf("%v -> \n", operationID)
	//		for _, route := range routes {
	//			fmt.Print(route)
	//		}
	//	}
	//	fmt.Println()
	//}
	fmt.Println(len(st.Lines))

	st.CrawlSchedules(1)
	//for k, v := range st.Schedules {
	//	fmt.Printf("%v,%v\n", k, v)
	//}
	fmt.Println(len(st.Schedules))

	st.CrawlVirtualTablesLines(crawler.Normal)
	//for _, vtStop := range st.VirtualTableStops {
	//	fmt.Println(vtStop)
	//}
	fmt.Println(len(st.VirtualTableStops))

	st.CrawlVirtualTablesStopsForTimes(100)
	//for k, v := range st.VirtualTableStopsTimes {
	//	fmt.Printf("%v -> %v\n", k, v)
	//}
	fmt.Println(len(st.VirtualTableStopsTimes))

}
