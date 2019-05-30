package workerpool 

import (
	"math/rand"
)

type Request struct {
	job 		func() int // the function to perform
	result 		chan int   // the channel to return the result
}


// Todo ..
func performJob() int {
	time.Sleep((time.Duration(rand.Intn(4)) * time.Second) + time.Second)
	return 1
}

func requester(requests chan Request) {
	result := make(chan int)
	for {

		// sleep for a while
		time.Sleep((time.Duration(rand.Intn(4)) * time.Second) + time.Second)

		select {
		// request sent
		case request  <- Request{performJob, result}:

		// result came back	
		case <-result:
		}
	}
}