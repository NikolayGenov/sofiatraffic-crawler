package crawler

import "github.com/garyburd/redigo/redis"

type stopTimesResponse struct {
	stop  VirtualTableStop
	times string
}

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

func (v *vtStopCrawler) createAndStartWorkers() {
	for i := range v.workers {
		v.workers[i] = newWorker(i+1, v.workerQueue, v.stopTimesQueue, v.finishChan)
		v.workers[i].start()
	}
}

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
func (v *vtStopCrawler) enqueueStops() {
	go func() {
		for _, stop := range v.stops {
			v.stopsQueue <- stop
		}
	}()
}

func (v *vtStopCrawler) waitForAllStops() {
	conn := v.pool.Get()
	defer conn.Close()
	for c := 0; c < len(v.stops); {
		select {
		case stopResponse := <-v.stopTimesQueue:
			(*v.stopTimesMap)[stopResponse.stop] = stopResponse.times
			conn.Do("HSET", "stops", stopResponse.stop, stopResponse.times)
		case <-v.finishChan:
			c++

		}
	}
}

func (v *vtStopCrawler) stop() {
	for _, worker := range v.workers {
		worker.stop()
	}
	close(v.done)
}
