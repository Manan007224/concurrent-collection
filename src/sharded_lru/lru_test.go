package lru

import (
	"math/rand"
	_ "net/http/pprof"
	"testing"
)

const nshards = 10000

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
	l := New(WithCapacity(10))
	for i := 0; i < 10; i++ {
		l.Add(i, i)
	}

	// cap was set to 10, so new adds should keep Len at 10
	for i := 10; i < 20; i++ {
		l.Add(i, i)
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

func makeRand(n int) []int {
	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = rand.Int()
	}
	return l
}

func BenchmarkAdd(b *testing.B) {
	b.Run("mostly_new", func(b *testing.B) {
		l := New(WithShards(nshards))
		rands := makeRand(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Add(rands[i], i)
		}
	})

	b.Run("mostly_existing", func(b *testing.B) {
		l := New(WithShards(nshards))
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

func BenchmarkAddParallel(b *testing.B) {
	b.Run("mostly_new", func(b *testing.B) {
		l := New(WithShards(nshards))
		rands := makeRand(b.N)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			// randomize a bit because otherwise
			// the goroutines will end up doing many
			// adds on the same key, which is not what
			// we're measuring here.
			n := rand.Int()
			for pb.Next() {
				l.Add(rands[i]+n, i)
				i++
			}
		})
	})

	b.Run("mostly_existing", func(b *testing.B) {
		l := New(WithShards(nshards))
		rands := makeRand(b.N)
		for i := 0; i < b.N; i++ {
			l.Add(rands[i], i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				l.Add(rands[i], i)
				i++
			}
		})
	})
}

var res interface{}

func populated(ns []int) *LRU {
	l := New(WithShards(nshards))
	for i, n := range ns {
		l.Add(n, i)
	}
	return l
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