package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

const (
	SCHEDULES_MAIN_URL            = "http://schedules.sofiatraffic.bg/"
	TRAMS_SECTION_PREFIX          = "Трамвайни линии"
	TROLLEY_SECTION_PREFIX        = "Тролейбусни линии"
	URBAN_BUSES_SECTION_PREFIX    = "Градски автобусни линии"
	SUBURBAN_BUSES_SECTION_PREFIX = "Крайградски автобусни линии"
	SUBWAY_SECTION                = "Метро"
)

//var rxOk = regexp.MustCompile(`http://schedules\.sofiatraffic\.bg(.*)$`)

type LineNameAndURL struct {
	name string
	url  string
}

type Lines struct {
	trams         []LineNameAndURL
	trolleys      []LineNameAndURL
	buses         []LineNameAndURL
	suburbanBuses []LineNameAndURL
	subwayLines   []LineNameAndURL
}

type LineTypesCrawler struct {
	gocrawl.DefaultExtender
	Lines
}

func (lt *LineTypesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	doc.Find(".lines_section").Each(func(i int, s *goquery.Selection) {
		typeAndLinesList := strings.Split(strings.TrimSpace(s.Text()), ":")
		lineType := typeAndLinesList[0]
		lines := strings.Fields(typeAndLinesList[1])
		linesInfo := make([]LineNameAndURL, len(lines))
		s.Find("a").Each(func(i int, link *goquery.Selection) {
			url, ok := link.Attr("href")
			if !ok {
				fmt.Errorf("The selector does not have an attribute href: %v", link.Text())
			}
			linesInfo[i] = LineNameAndURL{name: link.Text(), url: url}
		})
		switch lineType {
		case TRAMS_SECTION_PREFIX:
			lt.trams = linesInfo
		case TROLLEY_SECTION_PREFIX:
			lt.trolleys = linesInfo
		case URBAN_BUSES_SECTION_PREFIX:
			lt.buses = linesInfo
		case SUBURBAN_BUSES_SECTION_PREFIX:
			lt.suburbanBuses = linesInfo
		}
	})
	doc.Find(".quicksearch").Filter("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == SUBWAY_SECTION {
			url, ok := s.Attr("href")
			if !ok {
				fmt.Errorf("The selector does not have an attribute href: %v", s.Text())
			}
			lt.subwayLines = []LineNameAndURL{{name: SUBWAY_SECTION, url: url}}
		}
	})

	return nil, false
}

func crawlForListOfActiveLines() Lines {
	lineTypesCrawler := new(LineTypesCrawler)
	opts := gocrawl.NewOptions(lineTypesCrawler)
	opts.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	opts.CrawlDelay = 100 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 1
	opts.SameHostOnly = true
	crawler := gocrawl.NewCrawlerWithOptions(opts)
	crawler.Run(SCHEDULES_MAIN_URL)

	return lineTypesCrawler.Lines

}

func main() {
	lines := crawlForListOfActiveLines()
	fmt.Println("Trams:")
	fmt.Println(lines.trams)
	fmt.Println("Trolleys")
	fmt.Println(lines.trolleys)
	fmt.Println("Urban bus lines")
	fmt.Println(lines.buses)
	fmt.Println("Suburban bus lines")
	fmt.Println(lines.suburbanBuses)
	fmt.Println("Subway")
	fmt.Println(lines.subwayLines)
}
