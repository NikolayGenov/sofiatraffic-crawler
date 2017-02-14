package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type lineCrawler struct {
	gocrawl.DefaultExtender
	lines *[]Line
}

func (l *lineCrawler) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	path := ctx.URL().Path
	if path == "/" ||
		strings.HasPrefix(path, "/tramway") ||
		strings.HasPrefix(path, "/trolleybus") ||
		strings.HasPrefix(path, "/autobus") {
		return true
	}

	return false
}

func (l *lineCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	path := ctx.URL().Path
	if ctx.URL().Path != "/" {
		parts := strings.Split(path[1:], "/")
		name := parts[1]
		transportation, err := convertToTransportation(parts[0])
		if err != nil || name == "" {
			panic(fmt.Errorf("Unknown transporation type or empty name, given: %v", parts))
		}
		line := Line{
			Name:                 name,
			Transportation:       transportation,
			Path:                 path,
			OperationIDMap:       getOperationsMap(doc),
			OperationIDRoutesMap: getOperationIDRoutesMap(doc)}
		//TODO potential race
		*l.lines = append(*l.lines, line)

	}
	return nil, true
}

func newLineCrawler(lines *[]Line) crawlable {
	lineCrawler := &lineCrawler{lines: lines}
	opts := gocrawl.NewOptions(lineCrawler)
	opts.UserAgent = user_agent
	opts.CrawlDelay = 0
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = true
	c := gocrawl.NewCrawlerWithOptions(opts)
	return c
}
