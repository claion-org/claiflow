package generic

func MapE[A, B any](a []A, mapper func(a A) (B, error)) (b []B, err error) {
	b = make([]B, len(a))
	for i := range a {
		b[i], err = mapper(a[i])
		if err != nil {
			break
		}
	}

	return b, err
}

func Map[A, B any](a []A, mapper func(a A) B) (b []B) {
	b = make([]B, len(a))
	for i := range a {
		b[i] = mapper(a[i])
	}

	return b
}

func Fold[A, B any](acc B, a []A, folder func(acc B, a A) B) (b B) {
	for i := range a {
		acc = folder(acc, a[i])
	}

	return acc
}
