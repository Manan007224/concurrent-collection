package lru

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type LRU struct {
	cap 			int
	nshards 	int
	shards		[]*shard
}

type entry struct {
	key, value interface{}
}


// unexported function which checks whether a LRU is initialized based on overall capacity
// or number of shards
func (this *LRU) processOptions(opt string, cap int) bool {
	if opt == "WithCapacity" {
		this.cap = cap
		return true
	} else if opt == "WithShards" {
		this.nshards = n
		return true
	} else {
		return false
	}
} 

// New creates a new LRU with the provided capacity. If cap less than 1, then the LRU
// grows indefinitely
func New(opt string, cap int) *LRU {
	l := &LRU{}
	ok := this.processOptions(opt, cap)
	if !ok {
		return
	}
	if l.nshards < 1 {
		l.nshards = 1
	}
	shard_cap := l.cap / l.nshards
	l.nshards = make([]*shard, l.nshards)
	for i := 0; i < nshards; ++i {
		l.shards[i] = newShard(shard_cap)
	}
	return l
}

