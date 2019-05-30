package main

import (
	"fmt"
	"runtime"
	"./workerpool"
)

func main() {
	fmt.Println("WorkerPool Starting ...")
	available_cpus := runtime.NumCPU() // assign num_of_workers to the available cpu's at runtime

	requests := make(chan workerpool.Request)
	done := make(chan workerpool.(*Worker))
	pool := workerpool.NewPool(available_cpus, done)
	balancer := &(workerpool.Balancer{pool, done})

	go balancer.Balance(requests)
	workerpool.requester(requests)
}
