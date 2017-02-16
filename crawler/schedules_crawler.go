package crawler

import (
	"net/http"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type schedulesCrawler struct {
	gocrawl.DefaultExtender
	schedules *map[ScheduleID]ScheduleTimes
}

func newSchedulesCrawler(schedules *map[ScheduleID]ScheduleTimes) *gocrawl.Crawler {
	schedulesCrawler := &schedulesCrawler{
		schedules: schedules}
	opts := gocrawl.NewOptions(schedulesCrawler)
	opts.UserAgent = userAgent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	return gocrawl.NewCrawlerWithOptions(opts)
}

func (s *schedulesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	scheduleID := convertToScheduleID(ctx.URL().Path)
	scheduleTimes := getScheduleTimes(doc)
	(*s.schedules)[scheduleID] = scheduleTimes
	return nil, false
}
