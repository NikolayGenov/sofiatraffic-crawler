package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getTextAndHref(s *goquery.Selection) (string, string) {
	text := strings.TrimSpace(s.Text())
	if text == "" {
		html, _ := s.Html()
		panic(fmt.Errorf("The given attribute MUST have some text, HTML: %v", html))
	}

	href, ok := s.Attr("href")
	if !ok {
		html, _ := s.Html()
		panic(fmt.Errorf("The given attribute should have href attribute, HTML: %v", html))
	}
	return text, href
}
