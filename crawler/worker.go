package crawler

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

//worker represents one entity that is doing some work
//which in the case is crawling for up-to date stop times
type worker struct {
	//id is used for tracking purpose only
	id int

	//The requests the worker executes
	stopsChan chan VirtualTableStop

	//workerQueue is a channel where the worker should register when it is ready to do work
	workerQueue chan chan VirtualTableStop

	//responseQueue is a channel where the worker passes to the function that does the work
	//in order to return the results of its work in this channel
	responseQueue chan stopTimesResponse

	//responseQueue is a channel where the worker passes to the function that does the work
	//in order to to notify the listeners of the channel that it has finished its work
	finishChan chan struct{}

	//done is a channel which is used to stop the worker from doing any more work
	//Note that it does not stop the worker from doing its current work
	done chan struct{}
}

//newWorker returns an initialized worker, which have not started working yet - not self registered
func newWorker(id int, workerQueue chan chan VirtualTableStop, responseQueue chan stopTimesResponse, finishChan chan struct{}) worker {
	return worker{
		id:            id,
		stopsChan:     make(chan VirtualTableStop),
		workerQueue:   workerQueue,
		responseQueue: responseQueue,
		finishChan:    finishChan,
		done:          make(chan struct{})}
}

//start is used to make the worker register for work and then listens for work
//in this case  - for stop requests
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

//crawlStop takes a stop to process and a channels where to send the times when its done,
//and another channel to notify that it has finished
//It makes a post request to a not so popular link on a Virtual Tables site
//Then it has to parse the returned HTML
//Tokenizer is used because it is the simplest and possibly one of the fastest ways
//of getting the information that is needed
//The information (if present) is located on the forth "b" tag  e.g: <b>21:56:00,22:23:57,22:57:11</b>
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
