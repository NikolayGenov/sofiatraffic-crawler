package crawler

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

const (
	schedules_main_url          = "http://schedules.sofiatraffic.bg/"
	trams_query_prefix          = "Трамвайни линии"
	trolley_query_prefix        = "Тролейбусни линии"
	urban_buses_query_prefix    = "Градски автобусни линии"
	suburban_buses_query_prefix = "Крайградски автобусни линии"
	subway_query                = "Метро"
	user_agent                  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
)

type LineNameAndURL struct {
	Name string
	URL  string
}
type TransportationType int

const (
	Tram TransportationType = iota
	Trolley
	Bus
	Suburban
)

type Lines struct {
	Trams         []LineNameAndURL
	Trolleys      []LineNameAndURL
	Buses         []LineNameAndURL
	SuburbanBuses []LineNameAndURL
}

type lineTypesCrawler struct {
	gocrawl.DefaultExtender
	Lines
}

func (lt *lineTypesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	doc.Find(".lines_section").Each(func(i int, s *goquery.Selection) {
		typeAndLinesList := strings.Split(strings.TrimSpace(s.Text()), ":")
		lineType := typeAndLinesList[0]
		lineNames := strings.Fields(typeAndLinesList[1])
		linesInfo := make([]LineNameAndURL, len(lineNames))
		s.Find("a").Each(func(i int, link *goquery.Selection) {
			if url, ok := link.Attr("href"); ok {
				linesInfo[i] = LineNameAndURL{Name: link.Text(), URL: url}
			}
		})
		//FIXME - replace it with direct compare
		switch lineType {
		case trams_query_prefix:
			lt.Trams = linesInfo
		case trolley_query_prefix:
			lt.Trolleys = linesInfo
		case urban_buses_query_prefix:
			lt.Buses = linesInfo
		case suburban_buses_query_prefix:
			lt.SuburbanBuses = linesInfo
		}
	})

	return nil, false
}

func ActiveLines() Lines {
	lineTypesCrawler := new(lineTypesCrawler)
	opts := gocrawl.NewOptions(lineTypesCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 100 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.MaxVisits = 1
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(schedules_main_url)

	return lineTypesCrawler.Lines
}
