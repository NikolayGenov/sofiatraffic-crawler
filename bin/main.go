package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"../crawler"
	"github.com/garyburd/redigo/redis"
)

func loadIDLines(conn redis.Conn, paths []string) (lines []crawler.Line) {
	query := make([]string, len(paths))

	for i, path := range paths {
		query[i] = fmt.Sprintf("line%v", path)
	}

	serializedLines, _ := redis.Bytes(conn.Do("MGET", strings.Join(query, " ")))
	json.Unmarshal(serializedLines, &lines)
	return
}

func saveLines(conn redis.Conn, lines []crawler.Line) {
	for _, line := range lines {
		serialized, _ := json.Marshal(line)
		conn.Do("SET", fmt.Sprintf("line%v", line.Path), serialized)
	}
}

func loadAllLines(conn redis.Conn) (lines []crawler.Line) {
	serializedLines, _ := redis.Bytes(conn.Do("GET", "allLines"))
	json.Unmarshal(serializedLines, &lines)
	return
}

func crawlAllLines(conn redis.Conn, st *crawler.SofiaTrafficCrawler) []crawler.Line {
	st.CrawlLines()
	lines := st.Lines
	serialized, _ := json.Marshal(lines)
	conn.Do("SET", "allLines", serialized)

	return lines
}

func main() {
	st := crawler.NewSofiaTrafficCrawler()
	conn, _ := redis.Dial("tcp", ":6379")
	//start := time.Now()

	/* Load or crawl lines */
	//lines := crawlAllLines(conn)
	lines := loadAllLines(conn)
	st.Lines = lines[:1]
	//fmt.Println(len(lines))
	//elapsed := time.Since(start)
	//fmt.Printf("Took %s\n", elapsed)

	for _, l := range st.Lines {
		fmt.Println(l)
	}

	st.CrawlSchedules()

	for id, times := range st.Schedules {
		fmt.Printf("%v - > %v\n", id, times)
	}
	//seeds := make([]string, 0)
	//for _, line := range savedLines[:1] {
	//	scheduleIDs := line.ScheduleIDs()
	//	for _, id := range scheduleIDs {
	//		seeds = append(seeds, fmt.Sprintf("%v/%v", schedules_times_basic_url, id))
	//	}
	//}
	//schedules := crawler.CrawlSchedules(seeds)
	//fmt.Println(len(schedules))
	//fmt.Println(len(seeds))
	//elapsed2 := time.Since(start)
	//fmt.Printf("All Took %s\n", elapsed2)
}
