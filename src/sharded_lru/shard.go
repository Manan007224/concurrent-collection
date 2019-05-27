package lru

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type shard struct {
	cap 				int 												// The max no of items LRU can hold
	len 				int32
	cache 			map[interface{}]*list.Element // The cache for our items 
	evictList  *list.List 										// The acutal list holding our data
	sync.Mutex																// Protects the cache and evictList
}

func newShard(cap int) *shard {
	s := &shard{
		cap: 				 	cap,
		evictList:   	list.New(),
		cache: 			 	make(map[interface{}]*list.Element, cap+1),
	}
	return s
}

// Len returns the number of items currently in the LRU
func (s *shard) Len() int { return int(atomic.LoadInt32(&s.len)) }

// add will insert a new keyval pair to the shard
func (s *shard) add(k, v interface{}) {
	s.Lock()
	defer s.Unlock()

	// first let's see if we already have this key
	if le, ok := s.cache[k]; ok {
		// update the entry and move it to the front
		le.Value.(*entry).val = v
		s.evictList.MoveToFront(le)
		return
	}
	s.cache[k] = s.l.PushFront(&entry{key: k, val: v})
	atomic.AddInt32(&s.len, 1)

	if s.cap > 0 && s.Len() > s.cap {
		s.removeOldest()
	}
	return
}

// front will return the element at the front of the queue without modifying
// it in anyway
func (s *shard) front() (key, val interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.Len() == 0 {
		return nil, nil
	}

	le := s.evictList.Front()
	return le.Value.(*entry).key, le.Value.(*entry).val
}

// get will try to retrieve a value from the given key. The second return is
// true if the key was found.
func (s *shard) get(key interface{}) (value interface{}, ok bool) {
	s.Lock()
	defer s.Unlock()

	if le, found := s.cache[key]; found {
		s.l.MoveToFront(le)
		return le.Value.(*entry).val, true
	}
	return nil, false
}

func (s *shard) removeOldest() (key, val interface{}) {
	le := s.evictList.Back()
	if le == nil {
		return
	}
	return s.removeElement(le)
}

func (s *shard) removeElement(le *list.Element) (key, val interface{}) {
	e := le.Value.(*entry)
	s.evictList.Remove(le)
	delete(s.cache, e.key)
	atomic.AddInt32(&s.len, -1)
	return e.key, e.val
}

// removeKey will remove the given key from the LRU
func (s *shard) removeKey(key interface{}) {
	s.Lock()
	defer s.Unlock()

	le, ok := s.cache[key]
	if !ok {
		return
	}
	s.removeElement(le)
}