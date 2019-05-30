package pool

import (
	"fmt"
)

type Request struct {
	job 		func() int // the function to perform
	result 		chan int   // the channel to return the result
}