package crawler

var (
	workQueue          = make(chan workRequest, number_of_workers*2)
	workResponseQueues = make(chan workResponse, number_of_workers*2)
	finishedQueue      = make(chan struct{})
)

type workRequest struct {
	stop VirtualTableStop
}

type workerQueue chan chan workRequest

type workResponse struct {
	stop  VirtualTableStop
	times string
}

func startDispatcher(done chan struct{}, workerQueue workerQueue) {
	go func() {
		for {
			select {
			case work := <-workQueue:
				go func() {
					worker := <-workerQueue

					worker <- work
				}()
			case <-done:
				return
			}
		}
	}()
}

func stopWorkers(workers []worker) {
	for _, worker := range workers {
		worker.stop()
	}
}

func createAndStartWorkers(n int, workerQueue workerQueue) []worker {
	workers := make([]worker, n)
	for i := 0; i < n; i++ {
		workers[i] = newWorker(i+1, workerQueue)
		workers[i].start()
	}
	return workers

}
