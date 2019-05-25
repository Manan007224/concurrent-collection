// Implementation of LRU Cache which supports O(1) insertion and O(1) removal

package lru

import (
	"container/list"
)

type LRUCache struct {
	cap int // The max no of items LRUCache can hold
	cache map[inteface{}]*list.Element 	// The cache for our items 
	evict_list *list.List // The acutal list holding our data
}


// An unexported field which we actually store in our cache
type entry struct {
	key int
	value int
}

func New(cap int) *LRUCache {	
	lru_cache := &LRUCache {
		cap: cap
		evict_list: list.New()
		cache: make(map[interface{}]*list.Element, cap+1)
	}
	return lru_cache
} 

// Used to automatically initialize cache without the New method for eg:
// var L LRUCache
// L.Add("a", 5)

func (this* LRUCache) lazyInit {
	if this.evict_list == nil {
		this.evict_list = list.New()
		cache = make(map[interface{}]*list.Element, this.cap+1)
	}
} 

func (this *LRUCache) Add(k, v interface{}) {
	this.lazyInit()

	// If the item already exists
	if ent, ok := this.cache[k]; ok {
		this.evict_list.MoveToFront(ent)
		ent.Value.(*entry).value = v
		return
	}

	// If the capacity is full
	// Get the element which was least recently used from the evict_list
	if this.evict_list.Len() == this.cap {
		ent := this.evict_list.Back()
		k_v := ent.Value.(*entry)
		this.evict_list.Remove(ent)
		delete(this.cache, kv.key)
	}

	item := this.evict_list.PushFront(&entry{key, value})
	this.cache[key] = item
}


func (this *LRUCache) Get(k interface{}) {
	this.lazyInit()

	// Move the item at the head of the evict_list
	if ent, ok := this.cache[k]; ok {
		this.evict_list.MoveToFront(ent)
		return ent.Value.(*entry).value
	} else {
		return -1
	}
}

func (this *LRUCache) Remove(k, v interface{}) {
	this.lazyInit()

	ent, ok := this.cache[key]; ok {
		k_v := ent.Value.(*entry)
		this.evict_list.Remove(k_v)
		delete(this.cache, k_v.key)
		return ent.key, ent.val
	}
}