package main

import (
	"fmt"
	"container/heap"
)

type Pool []*Worker

const defaultSize int32 = 30

// create a new pool
func New(workers int, done chan *Request) *Pool {
	p := &Pool{}
	for w := 0; w < workers; ++w {
		requests := make(chan Request, defaultSize)
		worker := {requests, 0, i}
		go worker.Work(done)
		pool = append(pool, &worker)
	}
}

func (p *Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Len () int {
	return len(p)
}

func (p *Pool) swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i 
	p[j].index = j
}

func (p *Pool) Push(w interface{}) {
	worker := w.(*Worker)
	worker.index = p.Len()
	*p = append(*p,  worker)
}

func (p *Pool) Pop(w interface{}) {
	old := *(p)
	n := len(old)
	item := old[n-1]
	item.index = -1
	*(p) = old[0 : n-1]
	return item
}

