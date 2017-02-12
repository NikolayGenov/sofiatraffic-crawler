package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../crawler"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/redigo/redis"
)

const SCHEDULE_URL = "http://schedules.sofiatraffic.bg/"
const schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"

type x struct {
	gocrawl.DefaultExtender
	conn redis.Conn
}

func (lt *x) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	url := ctx.URL().Path[1:]
	html, _ := doc.Html()
	lt.conn.Do("SET", "wednesday/"+url, html)
	return nil, false
}

func allData(lines crawler.LinesBasicInfo) []string {
	seeds := make([]string, 0)
	for _, l := range lines {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	return seeds
}
func download(conn redis.Conn) {
	//lines := crawler.GetLineBasicInfo()
	//fmt.Println(lines)
	//lineInfo := lines[0]
	//line := crawler.CrawlLine(lineInfo)
	//fmt.Println(line)
	//fmt.Println(line.LinksToCrawl(schedules_times_basic_url))
	//serialized, _ := json.Marshal(lines)
	//conn.Do("SET", "lines", serialized)
	//
	//lineTypesCrawler := new(x)
	//lineTypesCrawler.conn = conn
	//opts := gocrawl.NewOptions(lineTypesCrawler)
	//opts.CrawlDelay = 1 * time.Millisecond
	//opts.LogFlags = gocrawl.LogEnqueued
	//
	//opts.SameHostOnly = true
	//c := gocrawl.NewCrawlerWithOptions(opts)
	//
	////Trams
	//
	//seeds := allData(&lines)
	//c.Run(seeds)

}

func loadLines(conn redis.Conn) (lines crawler.LinesBasicInfo) {
	serializedLines, _ := redis.Bytes(conn.Do("GET", "lines"))
	json.Unmarshal(serializedLines, &lines)
	return
}

func main() {

	conn, _ := redis.Dial("tcp", ":6379")
	//download(conn)

	lines := loadLines(conn)
	fmt.Println(lines)

	allLinks := make([]string, 0)
	for _, l := range lines {
		tramHTML, _ := redis.String(conn.Do("GET", "wednesday/"+l.URL))

		r := strings.NewReader(tramHTML)
		line := l.HelperTestReaderVisit(r)
		links := line.LinksToCrawl("")
		fmt.Println(line)
		allLinks = append(allLinks, links...)

	}
	fmt.Println(len(allLinks))

	//l := lines[1]
	//l := crawler.LineBasicInfo{"119", "autobus/119", crawler.Bus}
	//html, _ := redis.String(conn.Do("GET", "wednesday/"+l.URL))
	//r := strings.NewReader(html)
	//fmt.Println(l.HelperTestReaderVisit(r))
}
