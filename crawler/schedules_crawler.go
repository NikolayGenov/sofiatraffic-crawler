package crawler

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

//schedulesCrawler is extension to gocrawl. It takes a reference to a map with keys ScheduleID
//and values ScheduleTimes
type schedulesCrawler struct {
	gocrawl.DefaultExtender
	schedules *map[ScheduleID]ScheduleTimes
}

//newSchedulesCrawler takes a reference to a map with keys ScheduleID
//and values Schedule and creates an initialized instance of gocrawl.Crawler
//with proper user agent, delay and log options set
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

//Visit crawls all seed URLs and for each URL it first have to regenerate the ScheduleID
//Then search for the section where the schedules are shown and saves them to the map as a value
//with key the corresponding ScheduleID
func (s *schedulesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	scheduleID := convertToScheduleID(ctx.URL().Path)
	scheduleTimes := getScheduleTimes(doc)
	(*s.schedules)[scheduleID] = scheduleTimes
	return nil, false
}

func getScheduleTimes(doc *goquery.Document) ScheduleTimes {
	scheduleTimes := ScheduleTimes{}
	doc.Find(".schedule_times tbody a").
		Each(func(i int, s *goquery.Selection) {
			scheduleTimes = append(scheduleTimes, strings.TrimSpace(s.Text()))
		})
	return scheduleTimes
}
