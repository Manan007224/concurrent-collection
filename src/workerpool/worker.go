package workerpool 

type Worker struct {
	requests 	chan int	// All the pending requests(work to do ..)
	pending 	int			// count of remaining tasks
	index 		int			// index in the heap
}

// Worker performs the work to be done
func (w *Worker) Work(done chan *Worker) {
	for {
		select {
		case req := <-w.requests:
			req.result = req.performJob()
			done <- w
		}
	}
}