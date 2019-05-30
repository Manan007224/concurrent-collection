package main 

import (
	"fmt"
)

func main() {
	fmt.Println("WorkerPool Starting ...")
	available_cpus := runtime.NumCPU() // assign num_of_workers to the available cpu's at runtime

	requests := make(chan Request) // work to be done
	done := make(chan *Worker)

	pool := New(available_cpus, done)
	balancer := &balancer{pool, done}

	go balancer.Balance(requests)
	requester(requests)
}