package main

import (
	"net/http"

	"../crawler"

	"time"

	"encoding/json"

	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/redigo/redis"
)

const SCHEDULE_URL = "http://schedules.sofiatraffic.bg/"

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

func allData(lines *crawler.Lines) []string {
	seeds := make([]string, 0)
	for _, l := range lines.Trams {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	for _, l := range lines.Trolleys {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	for _, l := range lines.Buses {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	for _, l := range lines.SuburbanBuses {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	for _, l := range lines.SubwayLines {
		seeds = append(seeds, SCHEDULE_URL+l.URL)
	}
	return seeds
}
func download(conn redis.Conn) {
	lines := crawler.ActiveLines()
	serialized, _ := json.Marshal(lines)
	conn.Do("SET", "lines", serialized)

	lineTypesCrawler := new(x)
	lineTypesCrawler.conn = conn
	opts := gocrawl.NewOptions(lineTypesCrawler)
	opts.CrawlDelay = 1 * time.Millisecond
	opts.LogFlags = gocrawl.LogEnqueued

	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)

	//Trams

	seeds := allData(&lines)
	c.Run(seeds)

}

func loadLines(conn redis.Conn) crawler.Lines {

	var lines crawler.Lines
	serilizedLines, _ := redis.Bytes(conn.Do("GET", "lines"))
	json.Unmarshal(serilizedLines, &lines)
	return lines
}
func main() {

	conn, _ := redis.Dial("tcp", ":6379")

	//download(conn)

	//lines := loadLines(conn)
	//fmt.Println(lines)
	//for _, l := range lines.Trams {
	//	tramHTML, _ := redis.String(conn.Do("GET", "sunday/"+l.URL))
	//	r := strings.NewReader(tramHTML)
	//	crawler.CrawlLine(l, r)
	//}

	l := crawler.LineNameAndURL{"1", "tramway/1"}
	tramHTML, _ := redis.String(conn.Do("GET", "wednesday/"+l.URL))
	r := strings.NewReader(tramHTML)
	crawler.CrawlLine(l, r)

}
