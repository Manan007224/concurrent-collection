package main

import (
	"fmt"
)

type Request struct {
	job 		func() int // the function to perform
	result 		chan int   // the channel to return the result
}


// Todo ..
func performJob() int {

}

func requester(requests chan Request) {
	result := make(chan int)
	for {
		select {
		case request  <- Request{performJob, result}:

		case <-result:
		}
	}
}