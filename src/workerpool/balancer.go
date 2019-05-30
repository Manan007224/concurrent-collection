// Balancer sends request to most lightly loaded worker
// Uses min-heap

package main

import (
	"fmt"
	"container/heap"
)

type Balancer struct {
	workerpool 	*Pool
	done 		chan int
}

func (b *Balancer) Balance (work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(request)
		case worker := <-b.done:
			b.complete(worker)
	}
}

func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Worker) // Grabbing the least loaded worker
	w.requests <- req 
	w.pending++
	heap.Push(&b.pool, w)
}

func (b *Balancer) complete(w *Worker) {
	w.pending--;
	heap.Fix(&b.pool, w.index)
}