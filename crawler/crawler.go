package crawler

import (
	"fmt"

	"github.com/PuerkitoBio/gocrawl"
)

const (
	schedules_main_url        = "http://schedules.sofiatraffic.bg/"
	user_agent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
)

type SofiaTrafficCrawler struct {
	Lines []Line
	Schedules
}

func (s *SofiaTrafficCrawler) CrawlLines() {
	lineCrawler := &lineCrawler{}
	opts := gocrawl.NewOptions(lineCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url)
	s.Lines = lineCrawler.lines

}

func (s *SofiaTrafficCrawler) CrawlSchedules() {
	links := buildSchedulesLinks(s.Lines)
	schedulesCrawler := &schedulesCrawler{}
	schedulesCrawler.Schedules = make(Schedules)
	opts := gocrawl.NewOptions(schedulesCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(links)
	s.Schedules = schedulesCrawler.Schedules
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
