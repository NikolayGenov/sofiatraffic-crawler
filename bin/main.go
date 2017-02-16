package main

import (
	"fmt"
	"time"

	"../crawler"
)

func main() {

	st, err := crawler.NewSofiaTrafficCrawler(":6379")
	if err != nil {
		panic(err)
	}
	start := time.Now()
	//st.CrawlLines()
	//st.CrawlSchedules(1)

	//st.CrawlVirtualTablesLines(crawler.Operation_Normal)
	//fmt.Println(len(st.VirtualTableStops))
	st.CrawlVirtualTablesStopsForTimes(100)

	elapsed := time.Since(start)
	for k, v := range st.VirtualTableStopsTimes {
		fmt.Printf("%v -> %v\n", k, v)
	}
	fmt.Printf("Map: %v\n", len(st.VirtualTableStopsTimes))
	fmt.Printf("Took %s\n", elapsed)

}
