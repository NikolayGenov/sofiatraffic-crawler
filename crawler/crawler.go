package crawler

import (
	"time"

	"github.com/PuerkitoBio/gocrawl"
)

const (
	schedules_main_url        = "http://schedules.sofiatraffic.bg/"
	user_agent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
)

func GetLineBasicInfo() LinesBasicInfo {

	lineBasicInfoCrawler := &lineBasicInfoCrawler{}
	opts := gocrawl.NewOptions(lineBasicInfoCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 100 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 1
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url)

	return lineBasicInfoCrawler.LinesBasicInfo
}

func CrawlLine(line LineBasicInfo) Line {

	lineCrawler := &lineCrawler{}
	lineCrawler.LineBasicInfo = line
	opts := gocrawl.NewOptions(lineCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 100 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 1
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url + line.URL)

	return lineCrawler.Line
}
