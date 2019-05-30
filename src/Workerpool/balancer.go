package Workerpool

import (
	"container/heap"
	"fmt"
)

// Balancer has a Pool of workers and a channel to pass
// workers through when they are finished a task
type Balancer struct {
	Pool 	*Pool
	Done 	chan *Worker
}

// Balance takes in a channel of requests and distrubutes them
func (b *Balancer) Balance(requests <-chan Request) {
	for {
		select {
		case request := <-requests:
			b.dispatch(request)
			fmt.Println(b.Pool)
		case worker := <-b.Done:
			b.complete(worker)
		}
	}
}

// dispatch distrubutes the requests
func (b *Balancer) dispatch(request Request) {
	w := heap.Pop(b.Pool).(*Worker)
	w.requests <- request
	w.pending += 1
	heap.Push(b.Pool, w)
}

// complete updates the worker Pool when a request is complete
func (b *Balancer) complete(worker *Worker) {
	worker.pending -= 1
	heap.Fix(b.Pool, worker.index)
}