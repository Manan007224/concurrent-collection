package Workerpool 

import (
	"container/heap"
)

type Pool []*Worker

const defaultSize int32 = 30

// create a new pool
func New(workers int, done chan *Worker) *Pool {
	var p Pool
	for w := 0; w < workers; w++ {
		requests := make(chan Request, defaultSize)
		worker := Worker{requests, 0, w}
		go worker.Work(done)
		p = append(p, &worker)
	}
	heap.Init(&p)
	return &p
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len () int {
	return len(p)
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i 
	p[j].index = j
}

func (p *Pool) Push(w interface{}) {
	worker := w.(*Worker)
	worker.index = p.Len()
	*p = append(*p,  worker)
}

func (p *Pool) Pop() interface{} {
	old := *(p)
	n := len(old)
	item := old[n-1]
	item.index = -1
	*(p) = old[0 : n-1]
	return item
}

