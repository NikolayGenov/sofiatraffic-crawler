package crawler

import "github.com/garyburd/redigo/redis"

//stopTimesResponse is only used as a type for channel to return information from
// times crawling
type stopTimesResponse struct {
	stop  VirtualTableStop
	times string
}

//vtStopCrawler is a custom made crawler just for the purpose of crawling
// real-time-sh data from one endpoint.
// The crawler makes only POST requests
type vtStopCrawler struct {
	stopsQueue     chan VirtualTableStop
	stopTimesQueue chan stopTimesResponse
	workerQueue    chan chan VirtualTableStop
	workers        []worker
	done           chan struct{}
	finishChan     chan struct{}
	stops          []VirtualTableStop
	pool           *redis.Pool
	stopTimesMap   *map[VirtualTableStop]string
}

//newVTStopCrawler a new initializer vtStopCrawler crawler with some buffered channels
func newVTStopCrawler(stops []VirtualTableStop, stopTimesMap *map[VirtualTableStop]string, pool *redis.Pool) *vtStopCrawler {
	return &vtStopCrawler{
		stopsQueue:     make(chan VirtualTableStop, numberOfWorkers*2),
		stopTimesQueue: make(chan stopTimesResponse, numberOfWorkers*2),
		finishChan:     make(chan struct{}),
		workerQueue:    make(chan chan VirtualTableStop),
		workers:        make([]worker, numberOfWorkers),
		done:           make(chan struct{}),
		pool:           pool,
		stops:          stops,
		stopTimesMap:   stopTimesMap}

}

//createAndStartWorkers crates a new worker and place it to its position
// and only then starts it
func (v *vtStopCrawler) createAndStartWorkers() {
	for i := range v.workers {
		v.workers[i] = newWorker(i+1, v.workerQueue, v.stopTimesQueue, v.finishChan)
		v.workers[i].start()
	}
}

//startDispatcher is taking work from the work queue (stop in this example)
// then it finds a worker and gives that work to it by placing it in its channel
func (v *vtStopCrawler) startDispatcher() {
	go func() {
		for {
			select {
			case stop := <-v.stopsQueue:
				go func() {
					worker := <-v.workerQueue
					worker <- stop
				}()
			case <-v.done:
				return
			}
		}
	}()
}

//enqueueStops takes all channels and puts them on the work queue
// but it does it in a go routine so that it doesn't block anything else
func (v *vtStopCrawler) enqueueStops() {
	go func() {
		for _, stop := range v.stops {
			v.stopsQueue <- stop
		}
	}()
}

//waitForAllStops is waiting for all the stops to finish their work regardless of the result.
// For those that do return some useful result - it saves it both to a map with stop as key
// and sends to redis hash set with key SofiaTraffic/stops
func (v *vtStopCrawler) waitForAllStops() {
	conn := v.pool.Get()
	defer conn.Close()
	for c := 0; c < len(v.stops); {
		select {
		case stopResponse := <-v.stopTimesQueue:
			(*v.stopTimesMap)[stopResponse.stop] = stopResponse.times
			conn.Do("HSET", "SofiaTraffic/stops", stopResponse.stop, stopResponse.times)
		case <-v.finishChan:
			c++
		}
	}
}

//stop sends a stop signal on all the workers and then closes dispatcher.
func (v *vtStopCrawler) stop() {
	for _, worker := range v.workers {
		worker.stop()
	}
	close(v.done)
}
