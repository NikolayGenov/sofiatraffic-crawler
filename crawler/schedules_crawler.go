package crawler

import (
	"net/http"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type schedulesCrawler struct {
	gocrawl.DefaultExtender
	Schedules
}

func newSchedulesCrawler(schedules Schedules) runStopCapable {
	schedulesCrawler := &schedulesCrawler{
		Schedules: schedules}
	opts := gocrawl.NewOptions(schedulesCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	return gocrawl.NewCrawlerWithOptions(opts)
}

func (s *schedulesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	scheduleID := convertToScheduleID(ctx.URL().Path)
	scheduleTimes := getScheduleTimes(doc)
	s.Schedules[scheduleID] = scheduleTimes
	return nil, false
}
