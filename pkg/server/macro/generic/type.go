package generic

func Value[A any](p *A) A {
	if p == nil {
		var a A
		return a
	}

	return *p
}

func Pointer[A any](a A) *A {
	return &a
}

func IsNil[A any](p *A) bool {
	return p == nil
}

func Left[A, B any](a A, _ B) A {
	return a
}

func Right[A, B any](_ A, b B) B {
	return b
}
