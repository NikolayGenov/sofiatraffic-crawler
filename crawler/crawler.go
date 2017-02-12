package crawler

import (
	"time"

	"github.com/PuerkitoBio/gocrawl"
)

const (
	schedules_main_url = "http://schedules.sofiatraffic.bg/"
	user_agent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
)

func ActiveLines() LinesBasicInfo {

	lineTypesCrawler := &lineTypesCrawler{}
	opts := gocrawl.NewOptions(lineTypesCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 100 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 1
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url)

	return lineTypesCrawler.lines
}
