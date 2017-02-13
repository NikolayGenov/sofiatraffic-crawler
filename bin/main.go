package main

import (
	"fmt"

	"time"

	"../crawler"
	"github.com/garyburd/redigo/redis"
)

func main() {
	st := crawler.NewSofiaTrafficCrawler()

	conn, _ := redis.Dial("tcp", ":6379")
	start := time.Now()

	//st.CrawlLines()
	//st.SaveLines(conn)
	st.LoadLines(conn)
	fmt.Println(len(st.Lines))
	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)

	st.CrawlSchedules()
	st.SaveSchedules(conn)
	//st.LoadSchedules(conn)

	fmt.Println(len(st.Schedules))
	elapsed2 := time.Since(start)
	fmt.Printf("All Took %s\n", elapsed2)
}
