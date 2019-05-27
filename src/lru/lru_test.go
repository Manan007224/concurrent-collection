package lru

import (
	"math/rand"
	_ "net/http/pprof"
	"testing"
)

func TestEmptyValue(t *testing.T) {
	var l LRU
	for i := 0; i < 10; i++ {
		l.Add(i, i)
		if l.Len() != i+1 {
			t.Errorf("got len %d, want %d", l.Len(), i+1)
		}
	}
}

func TestNewCap(t *testing.T) {
	l := New(10)
	for i := 0; i < 10; i++ {
		l.Add(i, i)
		if l.Len() != i+1 {
			t.Errorf("got len %d, want %d", l.Len(), i+1)
		}
	}

	// cap was set to 10, so new adds should keep Len at 10
	for i := 10; i < 20; i++ {
		l.Add(i, i)
		if l.Len() != 10 {
			t.Errorf("got len %d, want 10", l.Len())
		}
	}
}

func TestAdd(t *testing.T) {
	var l LRU
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)
	if l.Len() != 3 {
		t.Errorf("got len %d, want 3", l.Len())
	}

	// readd an existing key
	l.Add(2, 10)
	if l.Len() != 3 {
		t.Errorf("got len %d, want 3", l.Len())
	}
	// also the value of 2 should be different now:
	if v, ok := l.Get(2); !ok {
		t.Error("could not get key 2")
	} else if v != 10 {
		t.Errorf("got %d, want 10", v)
	}

}

func TestRemoveOldest(t *testing.T) {
	l := New(3)
	for i := 0; i < 4; i++ {
		l.Add(i, i)
		if v, ok := l.Get(i); !ok {
			t.Errorf("could not get key %v", i)
		} else if v != i {
			t.Errorf("value of key %v is %v, want %v", i, v, i)
		}
	}
	// at this point, key 0 should have been pruned
	if _, ok := l.Get(0); ok {
		t.Error("key 0 should have been prune")
	}
}

func TestGetLatest(t *testing.T) {
	var l LRU
	// peek at an empty list
	k, v := l.GetLatest()
	if k != nil || v != nil {
		t.Errorf("GetLatest found something in empty LRU")
	}
	for i := 0; i < 10; i++ {
		l.Add(i, i)
	}
	if k, _ := l.GetLatest(); k != 9 {
		t.Errorf("front is %v, should be 9", k)
	}
	// make 5 hot but accessing it
	l.Get(5)
	// now 5 should be at the front
	if k, _ := l.GetLatest(); k != 5 {
		t.Errorf("front is %v, should be 5", k)
	}
}

func makeRand(n int) []int {
	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = rand.Int()
	}
	return l
}


func BenchmarkAdd(b *testing.B) {
	var l LRU
	b.Run("mostly_new", func(b *testing.B) {
		rands := makeRand(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Add(rands[i], i)
		}
	})

	b.Run("mostly_existing", func(b *testing.B) {
		rands := makeRand(b.N)
		for i := 0; i < b.N; i++ {
			l.Add(rands[i], i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Add(rands[i], i)
			i++
		}
	})
}

var res interface{}

func populated(ns []int) *LRU {
	var l LRU
	for i, n := range ns {
		l.Add(n, i)
	}
	return &l
}

func BenchmarkGet(b *testing.B) {
	b.Run("mostly_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res, _ = l.Get(r[i])
		}
	})

	b.Run("mostly_not_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res, _ = l.Get(i)
		}
	})
}

func BenchmarkGetParallel(b *testing.B) {
	b.Run("mostly_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				res, _ = l.Get(r[i]) // will find everything since we've added those numbers above
				i++
			}
		})
	})
	b.Run("mostly_not_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				res, _ = l.Get(i) // will rarely find a value
				i++
			}
		})
	})
}

func BenchmarkRemove(b *testing.B) {
	b.Run("mostly_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Remove(r[i]) // will find everything since we've added those numbers above
		}
	})

	b.Run("mostly_not_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Remove(i) // will rarely find a value
		}
	})
}

func BenchmarkRemoveParallel(b *testing.B) {
	b.Run("mostly_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				l.Remove(r[i]) // will find everything since we've added those numbers above
				i++
			}
		})
	})
	b.Run("mostly_not_found", func(b *testing.B) {
		r := makeRand(b.N)
		l := populated(r)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				l.Remove(i) // will rarely find a value
				i++
			}
		})
	})
}

func TestTraverse(t *testing.T) {
	var l LRU
	for i := 0; i < 1000; i++ {
		l.Add(i, i)
	}

	c := 0
	l.Traverse(func(key, val interface{}) bool {
		c++
		return true
	})
	if c != l.Len() {
		t.Errorf("c is %d, want %d", c, l.Len())
	}

	c = 0
	l.Traverse(func(key, val interface{}) bool {
		c++
		// stop traverse when c is 10
		return c != 10
	})
	if c != 10 {
		t.Errorf("c is %d, want 10", c)
	}
}