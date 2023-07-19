package generic

import "fmt"

type Lazy[A, B any] func(A) (B, error)

func Wrap[A, B any](fn func(A) B) Lazy[A, B] {
	return func(a A) (B, error) {
		return fn(a), nil
	}
}

func Placeholder[A any]() func(A) (A, error) {
	return func(a A) (A, error) {
		return a, nil
	}
}

func Compose[A, B, C any](lazyAB Lazy[A, B], lazyBC Lazy[B, C]) Lazy[A, C] {
	return func(a A) (c C, err error) {
		b, err := lazyAB(a)
		if err != nil {
			return c, err
		}

		return lazyBC(b)
	}
}

type aa = struct {
	string
	int
}

func AA(ss aa) {
	fmt.Print(ss.string, ss.int)
}
