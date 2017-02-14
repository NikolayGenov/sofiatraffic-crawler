package main

import (
	"fmt"

	"../crawler"
	"github.com/garyburd/redigo/redis"
)

func main() {
	st := crawler.NewSofiaTrafficCrawler()

	conn, _ := redis.Dial("tcp", ":6379")
	//start := time.Now()
	//elapsed := time.Since(start)
	//fmt.Printf("Took %s\n", elapsed)

	//st.CrawlLines()
	//st.SaveLines(conn)
	st.LoadLines(conn)
	fmt.Println(len(st.Lines))

	st.Lines = st.Lines[:1]
	operation := crawler.Operation_Normal
	st.CrawlVTLines(operation)
	fmt.Println(len(st.VirtualTableStops))
	for _, line := range st.Lines {
		id := line.OperationIDMap[operation]
		routes := line.OperationIDRoutesMap[id]
		for _, r := range routes {

			//fmt.Println(len(r.Stops))
			for _, s := range r.Stops {
				fmt.Println(s)
			}

		}
	}

	//st.CrawlSchedules()
	//st.SaveSchedules(conn)
	//st.LoadSchedules(conn)
	//fmt.Println(len(st.Schedules))
	//elapsed2 := time.Since(start)
	//fmt.Printf("All Took %s\n", elapsed2)
}
