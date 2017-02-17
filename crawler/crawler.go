package crawler

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	//Main source of information about Lines, Routes, Stops, and all kinds of IDs.
	schedulesMainURL = "http://schedules.sofiatraffic.bg/"

	//URL to an internal endpoint that provided HTML only one route, with schedule times for one stop.
	schedulesTimesBasicURL = "http://schedules.sofiatraffic.bg/server/html/schedule_load"

	//URL for a page that takes a POST form and returns HTML with real-time-sh times for each stop.
	virtualTableStopRealTimeURL = "http://m.sofiatraffic.bg/schedules/vehicle-vt"

	//Number of workers used for real-time-sh times crawling of Virtual Tables.
	numberOfWorkers = 30

	//Some sample user agent that is used by the polite crawler.
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
)

//SofiaTrafficCrawler struct keep all useful data that is extracted during
// different crawls.
type SofiaTrafficCrawler struct {
	//Internal redis connection pool for persistence
	redisPool *redis.Pool

	//List of active lines that were found during crawling
	Lines []Line

	//Map with keys unique string of type {OperationID}/{DirectionID}/{StopSign}
	Schedules map[ScheduleID]ScheduleTimes

	//List of active stops found on virtual tables site during crawling
	VirtualTableStops []VirtualTableStop

	//Map between those found stops and a string of comma separated times of arrival of the next vehicle
	VirtualTableStopsTimes map[VirtualTableStop]string
}

//NewSofiaTrafficCrawler creates an initialized NewSofiaTrafficCrawler struct that all crawler functions use.
// It takes already created pool of redis connections that it uses for persistence
// The data for is accessible trough the structure of the crawler.
func NewSofiaTrafficCrawler(redisPool *redis.Pool) *SofiaTrafficCrawler {
	return &SofiaTrafficCrawler{
		Lines:                  make([]Line, 0),
		Schedules:              make(map[ScheduleID]ScheduleTimes),
		VirtualTableStops:      make([]VirtualTableStop, 0),
		VirtualTableStopsTimes: make(map[VirtualTableStop]string),
		redisPool:              redisPool}
}

//CrawlLines starts a new crawl from schedules.sofiatraffic.bg as seed link and search for all links that match all
// transportation groups of links. Then for each found link, it parses the useful information and puts it
// into Lines variable on the SofiaTrafficCrawler struct.
// In the end it saves that information in redis
func (s *SofiaTrafficCrawler) CrawlLines() {
	lineCrawler := newLinesCrawler(&s.Lines)
	lineCrawler.Run(schedulesMainURL)
	s.saveLines()
}

//CrawlSchedules starts a new crawl by first building all the needed links from Lines.
// If it is an empty list - it loads it if it can from redis.
// The pages it crawls are from direct link from which gives only the schedules for one stop id.
// When crawling it saves the information corresponding to ScheduleID - which is list of time of day (24 hours)
// to a map which in the end saves to redis. It takes an int as a forNumberOfLines parameter which says
// how many of the found lines you want to crawl. If forNumberOfLines is 0, it crawls all the lines for schedule information
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
	links := buildSchedulesSeeds(lines)
	schedulesCrawler := newSchedulesCrawler(&s.Schedules)
	schedulesCrawler.Run(links)
	s.saveSchedules()
}

//CrawlVirtualTablesLines starts a new crawl by using the existing data from Lines.
// If it is an empty list, it loads it if it can from redis.
// It then builds links for crawling each line page in Virtual tables site. Note that there is
// significant differences in the data between Virtual Tables and Schedules hosted by sofiatraffic.bg.
// Meaning that no routes or stops match and it uses Schedules Sofia traffic as main source and
// only matches similar things in Virtual Tables site. It tries to parse and find all available stops for each line.
// When a stop is found it keeps it in the list VirtualTableStops - all the found active stops.
// It also updates the non capital name of a stop.
// In the end it saves the found stops in redis.
func (s *SofiaTrafficCrawler) CrawlVirtualTablesLines(operation Operation) {
	if len(s.Lines) == 0 {
		s.loadLines()
	}
	links := buildVirtualTablesSeeds(s.Lines)
	vtCrawler := newVirtualTableLineCrawler(s.Lines, &s.VirtualTableStops, operation)
	vtCrawler.Run(links)
	s.saveVirtualTableStops()
}

//CrawlVirtualTablesStopsForTimes stats a new crawl by using VirtualTableStops
// If it is an empty list - it loads it if it can from redis.
// It uses simpler and faster crawler which visits only 1 type of page
// process it's HTML by looking for specific ordering and extract comma separated times string
// It takes as int parameter forNumberOfStops which says how many of the already loaded
// virtual stops to crawl. If the parameter is 0, then it crawls all the available stops
func (s *SofiaTrafficCrawler) CrawlVirtualTablesStopsForTimes(forNumberOfStops int) {
	if len(s.VirtualTableStops) == 0 {
		s.loadVirtualTableStops()
	}
	var stops []VirtualTableStop
	if forNumberOfStops != 0 && forNumberOfStops <= len(s.VirtualTableStops) {
		stops = s.VirtualTableStops[:forNumberOfStops]
	} else {
		stops = s.VirtualTableStops
	}
	vtStopsCrawler := newVTStopCrawler(stops, &s.VirtualTableStopsTimes, s.redisPool)
	vtStopsCrawler.createAndStartWorkers()
	vtStopsCrawler.startDispatcher()
	vtStopsCrawler.enqueueStops()
	vtStopsCrawler.waitForAllStops()
	vtStopsCrawler.stop()

}

//Create seed URLs in the format that Sofia Traffic's internal server requires.
// That format is represented here in ScheduleID.
// For each line there are plenty of potential schedule URLs that can be crawled.
// Returns a slice of all seed URLs.
func buildSchedulesSeeds(lines []Line) []string {
	scheduleLinks := make([]string, 0)
	for _, line := range lines {
		for _, id := range line.scheduleIDs() {
			scheduleLinks = append(scheduleLinks, fmt.Sprintf("%v/%v", schedulesTimesBasicURL, id))
		}
	}
	return scheduleLinks
}

//Create seed URLs in the format that Virtual Table line schedules requires
// Information from a Line is enough to create one seed url for virtual tables line page.
// Returns a slice of all seed URLs.
func buildVirtualTablesSeeds(lines []Line) []string {
	vtScheduleLinks := make([]string, 0)
	for _, line := range lines {
		vtScheduleLinks = append(vtScheduleLinks,
			fmt.Sprintf(
				"http://m.sofiatraffic.bg/schedules?tt=%v&ln=%v&s=Търсене",
				int(line.Transportation),
				line.Name))
	}
	return vtScheduleLinks
}
