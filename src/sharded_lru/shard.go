package lru

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type shard struct {
	cap 				int32 												// The max no of items LRU can hold
	len 				int32
	cache 			map[interface{}]*list.Element // The cache for our items 
	evictList  *list.List 										// The acutal list holding our data
	sync.Mutex																// Protects the cache and evictList
}

