package lru

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"strconv"
)

type LRU struct {
	cap 		int
	nshards 	int
	shards		[]*shard
}

type entry struct {
	key, value interface{}
}


// unexported function which checks whether a LRU is initialized based on overall capacity
// or number of shards
// func (this *LRU) processOptions(opt string, cap int) bool {
// 	if opt == "WithCapacity" {
// 		this.cap = cap
// 		return true
// 	} else if opt == "WithShards" {
// 		this.nshards = n
// 		return true
// 	} else {
// 		return false
// 	}
// } 

// New creates a new LRU with the provided capacity. If cap less than 1, then the LRU
// grows indefinitely
func New(opts ...Option) *LRU {
	l := &LRU{}
	for _, o := range opts {
		o.apply(l)
	}
	if l.nshards < 1 {
		l.nshards = 1
	}

	cap := l.cap / l.nshards
	l.shards = make([]*shard, l.nshards)
	for i := 0; i < l.nshards; i++ {
		l.shards[i] = newShard(cap)
	}
	return l
}

// this initializes some fields at first use. Helpful to
// allow us to use the empty value of LRU
func (this *LRU) lazyInit() {
	if this.shards == nil {
		this.nshards = 1
		this.shards = []*shard{newShard(this.cap)}
	}
}

func (this *LRU) Len() int {
	this.lazyInit()
	var len int
	for i := 0; i < this.nshards; i++ {
		len += this.shards[i].Len()
	}
	return len
}

// Add will insert a new keyval pair to the LRU
func (this *LRU) Add(k, v interface{}) {
	this.lazyInit()
	this.shard(k).add(k, v)
}

// PeekFront will return the element at the front of the queue without modifying
// it in anyway
func (this *LRU) PeekFront() (key, val interface{}) {
	this.lazyInit()
	return this.shard(1).front()
}

// Get will try to retrieve a value from the given key. The second return is
// true if the key was found.
func (this *LRU) Get(key interface{}) (value interface{}, ok bool) {
	return this.shard(key).get(key)
}

// Remove will remove the given key from the LRU
func (this *LRU) Remove(key interface{}) {
	this.shard(key).removeKey(key)
}

// TraverseFunc is the function called for each element when
// traversing an LRU
type TraverseFunc func(key, val interface{}) bool

// Traverse will call fn for each element in the LRU, from most recently used to
// least. If fn returns false, the traverse stops
func (this *LRU) Traverse(fn TraverseFunc) {
L:
	for _, s := range this.shards {
		le := s.evictList.Front()
		for {
			if le == nil {
				break L
			}

			e := le.Value.(*entry)
			if !fn(e.key, e.value) {
				break L
			}
			le = le.Next()
		}
	}
}

// TraverseReverse will call fn for each element in the LRU, from least recently used to
// most. If fn returns false, the traverse stops
func (this *LRU) TraverseReverse(fn TraverseFunc) {
L:
	for _, s := range this.shards {
		le := s.evictList.Back()
		for {
			if le == nil {
				break L
			}

			e := le.Value.(*entry)
			if !fn(e.key, e.value) {
				break L
			}
			le = le.Prev()
		}
	}
}



// useful Interfaces for our hash function
type stringer interface {
	String() string
}

type byter interface {
	Bytes() []byte
}

// Implementation of the Fowler–Noll–Vo hash function. Simple but very useful one.

func (this *LRU) shard(key interface{}) *shard {
	h := fnv.New32a() // used to hash a byte array

	// try to get a bytes representation of the key any way we can, in order
	// from fastest to slowest
	switch v := key.(type) {
	case []byte:
		h.Write(v)
	case byter:
		h.Write(v.Bytes())
	case string:
		h.Write([]byte(v))
	case stringer:
		h.Write([]byte(v.String()))
	case int:
		h.Write(intBytes(v))
	case *int:
		h.Write(intBytes(*v))
	case *bool, bool, []bool, *int8, int8, []int8, *uint8,
		uint8, *int16, int16, []int16, *uint16,
		uint16, []uint16, *int32, int32, []int32, *uint32, uint32, []uint32,
		*int64, int64, []int64, *uint64, uint64, []uint64:
		h.Write(toBytes(v))
	default:
		// If the user is using an unknown type of key. A bit slower than the above ones.
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(v)
		if err != nil {
			panic(fmt.Sprintf("could not encode type %T as bytes", key))
		}
		h.Write(buf.Bytes())
	}

	return this.shards[h.Sum32()&uint32(this.nshards-1)]
}

func toBytes(v interface{}) []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
		panic(fmt.Sprintf("could not encode %v as bytes: %v", v, err))
	}
	return buf.Bytes()
}

var il = strconv.IntSize / 8

// helper function to quickly turn an int into a byte slice
func intBytes(i int) []byte {
	b := make([]byte, il)
	b[0] = byte(i)
	b[1] = byte(i >> 8)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 24)
	if il == 8 {
		b[4] = byte(i >> 32)
		b[5] = byte(i >> 40)
		b[6] = byte(i >> 48)
		b[7] = byte(i >> 56)
	}
	return b
}