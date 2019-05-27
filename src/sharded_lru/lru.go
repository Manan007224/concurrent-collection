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
		// the user is using an unknown type as the key, so we're now grasping
		// at straws here. This will be at least an order of magnitude slower
		// then the options above.
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