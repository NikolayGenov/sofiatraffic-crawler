package crawler

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	schedules_main_url                            = "http://schedules.sofiatraffic.bg/"
	user_agent                                    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url                     = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
	virtual_tables_url_placeholder_url            = "http://m.sofiatraffic.bg/schedules?tt=%v&ln=%v&s=Търсене"
	virtual_table_stop_real_time_link             = "http://m.sofiatraffic.bg/schedules/vehicle-vt"
	number_of_workers                             = 30
	stop_times_position_on_page_relative_to_bolds = 4
)

type SofiaTrafficCrawler struct {
	redisPool *redis.Pool
	Lines     []Line
	Schedules
	VirtualTableStops      []VirtualTableStop
	VirtualTableStopsTimes map[VirtualTableStop]string
}

type runStopCapable interface {
	Run(seeds interface{}) error
	Stop()
}

func NewSofiaTrafficCrawler(redisAddress string) (*SofiaTrafficCrawler, error) {
	pool := newPool(redisAddress)
	c, err := pool.Dial()
	defer c.Close()
	return &SofiaTrafficCrawler{
		Lines:                  make([]Line, 0),
		Schedules:              make(map[ScheduleID]ScheduleTimes),
		VirtualTableStops:      make([]VirtualTableStop, 0),
		VirtualTableStopsTimes: make(map[VirtualTableStop]string),
		redisPool:              pool}, err
}

func (s *SofiaTrafficCrawler) CrawlLines() {
	lineCrawler := newLineCrawler(&s.Lines)
	lineCrawler.Run(schedules_main_url)
	s.saveLines()

}

func (s *SofiaTrafficCrawler) CrawlSchedules(forNumberOfLines int) {
	if len(s.Lines) == 0 {
		s.loadLines()
	}
	var lines []Line
	if forNumberOfLines != 0 && forNumberOfLines <= len(s.Lines) {
		lines = s.Lines[:forNumberOfLines]
	} else {
		lines = s.Lines
	}

	links := buildSchedulesLinks(lines)
	schedulesCrawler := newSchedulesCrawler(s.Schedules)
	schedulesCrawler.Run(links)
	s.saveSchedules()
}

func (s *SofiaTrafficCrawler) CrawlVirtualTablesLines(operation Operation) {
	if len(s.Lines) == 0 {
		s.loadLines()
	}
	links := buildVirtualTablesLinks(s.Lines)
	vtCrawler := newVirtualTableLineCrawler(s.Lines, &s.VirtualTableStops, operation)
	vtCrawler.Run(links)
	s.saveVirtualTableStops()
}

func (s *SofiaTrafficCrawler) CrawlVirtualTablesStopsForTimes(forNumberOfStops int) {
	if len(s.VirtualTableStops) == 0 {
		s.loadVirtualTableStops()
	}
	done := make(chan struct{})
	defer close(done)
	if len(s.VirtualTableStops) == 0 {
		s.loadVirtualTableStops()
	}
	var stops []VirtualTableStop

	if forNumberOfStops != 0 && forNumberOfStops <= len(s.VirtualTableStops) {
		stops = s.VirtualTableStops[:forNumberOfStops]
	} else {
		stops = s.VirtualTableStops
	}

	workerQueue := make(workerQueue)
	workers := createAndStartWorkers(number_of_workers, workerQueue)
	startDispatcher(done, workerQueue)

	go func() {
		for _, stop := range stops {
			workQueue <- workRequest{stop: stop}
		}
	}()
	conn := s.redisPool.Get()
	defer conn.Close()
	for c := 0; c < len(stops); {
		select {
		case stopResponse := <-workResponseQueues:
			s.VirtualTableStopsTimes[stopResponse.stop] = stopResponse.times
			conn.Do("HSET", "stops", stopResponse.stop, stopResponse.times)
		case <-finishedQueue:
			c++
		}
	}
	stopWorkers(workers)
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

func buildVirtualTablesLinks(lines []Line) []string {
	scheduleLinks := make([]string, 0)
	for _, line := range lines {
		scheduleLinks = append(scheduleLinks, fmt.Sprintf(virtual_tables_url_placeholder_url, int(line.Transportation), line.Name))
	}
	return scheduleLinks
}
