package crawler

import (
	"time"

	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

const (
	schedules_main_url        = "http://schedules.sofiatraffic.bg/"
	user_agent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
	schedules_times_basic_url = "http://schedules.sofiatraffic.bg/server/html/schedule_load"
)

type lineBasicInfoCrawler struct {
	gocrawl.DefaultExtender
	LinesBasicInfo
}

type lineCrawler struct {
	gocrawl.DefaultExtender
	Line
}

func (l *lineBasicInfoCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	l.LinesBasicInfo = getLinesBasicInfo(doc)
	return nil, false
}

func (l *lineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	l.OperationIDMap = getOperationsMap(doc)
	l.OperationIDRoutesMap = getOperationIDRoutesMap(doc)
	return nil, false
}

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

func (l *LineBasicInfo) HelperTestReaderVisit(r io.Reader) Line {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	line := Line{}
	line.LineBasicInfo = *l
	line.OperationIDMap = getOperationsMap(doc)
	line.OperationIDRoutesMap = getOperationIDRoutesMap(doc)

	return line
}
