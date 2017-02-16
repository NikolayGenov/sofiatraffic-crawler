package crawler

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type worker struct {
	id            int
	stopsChan     chan VirtualTableStop
	workerQueue   chan chan VirtualTableStop
	done          chan struct{}
	responseQueue chan stopTimesResponse
	finishChan    chan struct{}
}

func newWorker(id int, workerQueue chan chan VirtualTableStop, responseQueue chan stopTimesResponse, finishChan chan struct{}) worker {
	return worker{
		id:            id,
		stopsChan:     make(chan VirtualTableStop),
		workerQueue:   workerQueue,
		done:          make(chan struct{}),
		responseQueue: responseQueue,
		finishChan:    finishChan,
	}
}

func (w *worker) start() {
	go func() {
		for {
			//Self register
			w.workerQueue <- w.stopsChan

			select {
			case stop := <-w.stopsChan:
				crawlStop(stop, w.responseQueue, w.finishChan)

			case <-w.done:
				return
			}
		}
	}()
}

func (w worker) stop() {
	close(w.done)
}

func crawlStop(stop VirtualTableStop, responseQueue chan stopTimesResponse, finishChan chan struct{}) {
	resp, err := http.PostForm(virtualTableStopRealTimeURL,
		url.Values{
			"vt":   {stop.TransportationType},
			"rid":  {stop.RouteID},
			"lid":  {stop.LineID},
			"stop": {stop.StopID}})

	defer func() { finishChan <- struct{}{} }()

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
			if boldCounter == 4 {
				if html.EndTagToken != tokenizer.Next() {
					responseQueue <- stopTimesResponse{stop, string(tokenizer.Raw())}
					return
				}
				return
			}
		}
	}
}
