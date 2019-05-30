// Balancer sends request to most lightly loaded worker
// Uses min-heap

package pool

import (
	"fmt"
)

type Balancer struct {
	workerpool 	*Pool
	done 		chan int
}