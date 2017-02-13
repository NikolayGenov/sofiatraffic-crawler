package crawler

import (
	"net/http"

	"fmt"

	"strings"

	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

const (
	schedules_main_url        = "http://schedules.sofiatraffic.bg/"
	user_agent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
)

type lineCrawler struct {
	gocrawl.DefaultExtender
	lines []Line
}

type schedulesCrawler struct {
	gocrawl.DefaultExtender
	Schedules
}

func (l *lineCrawler) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	path := ctx.URL().Path
	if path == "/" ||
		strings.HasPrefix(path, "/tramway") ||
		strings.HasPrefix(path, "/trolleybus") ||
		strings.HasPrefix(path, "/autobus") {
		return true
	}

	return false
}

func (l *lineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	path := ctx.URL().Path
	if ctx.URL().Path != "/" {
		parts := strings.Split(path[1:], "/")
		name := parts[1]
		transportation, err := convertToTransportation(parts[0])
		if err != nil || name == "" {
			panic(fmt.Errorf("Unknown transporation type or empty name, given: %v", parts))
		}
		line := Line{
			Name:                 name,
			Transportation:       transportation,
			Path:                 path,
			OperationIDMap:       getOperationsMap(doc),
			OperationIDRoutesMap: getOperationIDRoutesMap(doc)}
		l.lines = append(l.lines, line)
	}
	return nil, true
}

func (s *schedulesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	scheduleID := convertToScheduleID(ctx.URL().Path)
	scheduleTimes := getScheduleTimes(doc)
	s.Schedules[scheduleID] = scheduleTimes
	return nil, false
}

func CrawlLines() []Line {
	lineCrawler := &lineCrawler{}
	opts := gocrawl.NewOptions(lineCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url)

	return lineCrawler.lines
}

func CrawlSchedules(lines []Line) Schedules {
	links := buildSchedulesLinks(lines)
	schedulesCrawler := &schedulesCrawler{}
	schedulesCrawler.Schedules = make(Schedules)
	opts := gocrawl.NewOptions(schedulesCrawler)
	opts.UserAgent = user_agent

	opts.CrawlDelay = 1 * time.Microsecond
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(links)

	return schedulesCrawler.Schedules
}

func buildSchedulesLinks(lines []Line) []string {
	scheduleLinks := make([]string, 0)
	for _, line := range lines {
		for _, id := range line.ScheduleIDs() {
			scheduleLinks = append(scheduleLinks, fmt.Sprintf("%v/%v", schedules_times_basic_url, id))
		}
	}
	return scheduleLinks
}
