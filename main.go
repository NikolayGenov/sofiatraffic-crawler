package main

import (
	"fmt"

	"time"

	"github.com/NikolayGenov/sofiatraffic-crawler/crawler"
	"github.com/garyburd/redigo/redis"
)

//newPool returns a new initialized redis pool for connections
// It takes address e.g ":6379" as a parameter and uses it in the Dial function
func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 360 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", address) }}
}

func main() {
	pool := newPool(":6379")
	st := crawler.NewSofiaTrafficCrawler(pool)

	st.CrawlLines()
	for _, line := range st.Lines {
		fmt.Printf("Line [%v %v]\nOperations:  %v\n", line.Transportation, line.Name, &line.OperationIDMap)
		for operationID, routes := range line.OperationIDRoutesMap {
			fmt.Printf("%v -> \n", operationID)
			for _, route := range routes {
				fmt.Print(route)
			}
		}
		fmt.Println()
	}
	fmt.Println(len(st.Lines))

	st.CrawlSchedules(1)
	for k, v := range st.Schedules {
		fmt.Printf("%v,%v\n", k, v)
	}
	fmt.Println(len(st.Schedules))

	st.CrawlVirtualTablesLines(crawler.Normal)
	//for _, vtStop := range st.VirtualTableStops {
	//	fmt.Println(vtStop)
	//}
	fmt.Println(len(st.VirtualTableStops))

	st.CrawlVirtualTablesStopsForTimes(100)
	for k, v := range st.VirtualTableStopsTimes {
		fmt.Printf("%v -> %v\n", k, v)
	}
	fmt.Println(len(st.VirtualTableStopsTimes))

}
