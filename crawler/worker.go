package crawler

type worker struct {
	id          int
	workChan    chan workRequest
	workerQueue chan chan workRequest
	stopChan    chan struct{}
}

func newWorker(id int, workerQueue chan chan workRequest) worker {
	return worker{
		id:          id,
		workChan:    make(chan workRequest),
		workerQueue: workerQueue,
		stopChan:    make(chan struct{}),
	}
}

func (w *worker) start() {
	go func() {
		for {
			//Self register
			w.workerQueue <- w.workChan

			select {
			case work := <-w.workChan:
				crawlStop(work.stop)

			case <-w.stopChan:
				return
			}
		}
	}()
}

func (w worker) stop() {
	close(w.stopChan)
}
