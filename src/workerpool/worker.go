package pool

import (
	"fmt"
)

type Worker struct {
	requests 	chan int // All the pending requests(work to do ..)
	pending 	int			 // count of remaining tasks
	index 		int			 // index in the heap
}

