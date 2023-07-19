package generic

func Append[A any](dest []A, e ...A) []A {
	if cap(dest) < len(dest)+len(e) {
		swap := make([]A, len(dest), 2*(len(dest)+len(e)))
		copy(swap, dest)
		dest = swap
	}

	return append(dest, e...)
}

func Foreach[A any](aa []A, iter func(i int, a A) bool) {
	for i := range aa {
		if !iter(i, aa[i]) {
			break
		}
	}
}
