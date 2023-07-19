package webhook

import "sync"

type Pool[T any] struct {
	pool *sync.Pool
}

func NewPool[T any](new_ func() *T) *Pool[T] {
	conv := func(new_ func() *T) func() any {
		return func() any { return new_() }
	}

	return &Pool[T]{pool: &sync.Pool{New: conv(new_)}}
}

func (pool Pool[T]) Get() *T {
	return pool.pool.Get().(*T)
}

func (pool Pool[T]) Put(x *T) {
	pool.pool.Put(x)
}
