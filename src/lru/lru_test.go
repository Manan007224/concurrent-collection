package lru

import (
	// "math/rand"
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