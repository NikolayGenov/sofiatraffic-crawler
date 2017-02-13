package crawler

import "fmt"

const (
	schedules_main_url        = "http://schedules.sofiatraffic.bg/"
	user_agent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
)

type SofiaTrafficCrawler struct {
	Lines []Line
	Schedules
}

type crawlable interface {
	Run(seeds interface{}) error
	Stop()
}

func NewSofiaTrafficCrawler() *SofiaTrafficCrawler {
	return &SofiaTrafficCrawler{
		Lines:     make([]Line, 0),
		Schedules: make(map[ScheduleID]ScheduleTimes)}
}

func (s *SofiaTrafficCrawler) CrawlLines() {
	lineCrawler := newLineCrawler(&s.Lines)
	lineCrawler.Run(schedules_main_url)
}

func (s *SofiaTrafficCrawler) CrawlSchedules() {
	links := buildSchedulesLinks(s.Lines)
	schedulesCrawler := newSchedulesCrawler(s.Schedules)
	schedulesCrawler.Run(links)
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
