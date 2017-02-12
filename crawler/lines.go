package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

var transportationTypeIdentifierMap = map[string]TransportationType{
	"tramway":    Tram,
	"trolleybus": Trolley,
	"autobus":    Bus}

type LineBasicInfo struct {
	Name string
	URL  string
	Type TransportationType
}

type LinesBasicInfo []LineBasicInfo

type lineTypesCrawler struct {
	gocrawl.DefaultExtender
	lines LinesBasicInfo
}

func (l *LinesBasicInfo) getLines(doc *goquery.Document) {
	doc.Find(".lines_section ul li a").
		Each(func(i int, link *goquery.Selection) {
			url, ok := link.Attr("href")
			if !ok {
				html, _ := link.Html()
				panic(fmt.Errorf("Link should have 'href' attribute in order to process it, HTML: %v", html))
			}
			lineTypeString := strings.Split(url, "/")[0]
			lineType, ok := transportationTypeIdentifierMap[lineTypeString]
			if !ok {
				panic(fmt.Errorf("Line MUST one of required types in order to be processed, given: %v", lineTypeString))
			}
			*l = append(*l, LineBasicInfo{Name: link.Text(), URL: url, Type: lineType})
		})
}

func (lt *lineTypesCrawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	lt.lines.getLines(doc)
	return nil, false
}
