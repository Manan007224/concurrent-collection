package lru

// Option configures the LRU
type Option interface {
	apply(*LRU)
}

type optionFn func(*LRU)

func (f optionFn) apply(l *LRU) {
	f(l)
}

// WithCapacity configures the LRU to have a maximum capacity
func WithCapacity(cap int) Option {
	return optionFn(func(l *LRU) {
		l.cap = cap
	})
}

// WithShards configures the LRU to use the specified number of shards
func WithShards(n int) Option {
	return optionFn(func(l *LRU) {
		l.nshards = n
	})
}