package main

import (
	"fmt"
	"runtime"
	"./Workerpool"
)

func main() {
	fmt.Println("Workerpool Starting ...")
	available_cpus := runtime.NumCPU() // assign num_of_workers to the available cpu's at runtime

	requests := make(chan Workerpool.Request)
	done := make(chan *Workerpool.Worker)
	pool := Workerpool.New(available_cpus, done)
	balancer := &Workerpool.Balancer{pool, done}

	go balancer.Balance(requests)
	Workerpool.Requester(requests)
}
