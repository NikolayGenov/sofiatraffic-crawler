package crawler

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func crawlStop(stop VirtualTableStop) {
	resp, err := http.PostForm(virtual_table_stop_real_time_link,
		url.Values{
			"vt":   {stop.TransportationType},
			"rid":  {stop.RouteID},
			"lid":  {stop.LineID},
			"stop": {stop.StopID}})

	defer func() {
		finishedQueue <- struct{}{}
	}()

	if err != nil {
		log.Printf("Failed to fetch stop [stop: %v rid: %v lid: %v vt: %v]", stop.StopID, stop.RouteID, stop.LineID, stop.TransportationType)
		return
	}

	body := resp.Body
	defer body.Close()

	boldCounter := 0
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			isBold := token.Data == "b"
			if !isBold {
				continue
			}

			boldCounter++
			if boldCounter == stop_times_position_on_page_relative_to_bolds {
				if html.EndTagToken != tokenizer.Next() {
					workResponseQueues <- workResponse{stop, string(tokenizer.Raw())}
					return
				}
				return
			}
		}
	}
}
