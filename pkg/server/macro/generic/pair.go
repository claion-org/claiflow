package generic

type Pair[A, B any] struct {
	A A
	B B
}

func (p Pair[A, B]) Deconstruct() (A, B) {
	return p.A, p.B
}

type Tuple3[A, B, C any] struct {
	A A
	B B
	C C
}

func (patr Tuple3[A, B, C]) Deconstruct() (A, B, C) {
	return patr.A, patr.B, patr.C
}

type Tuple4[A, B, C, D any] struct {
	A A
	B B
	C C
	D D
}

func (t Tuple4[A, B, C, D]) Deconstruct() (A, B, C, D) {
	return t.A, t.B, t.C, t.D
}

type Tuple5[A, B, C, D, E any] struct {
	A A
	B B
	C C
	D D
	E E
}

func (t Tuple5[A, B, C, D, E]) Deconstruct() (A, B, C, D, E) {
	return t.A, t.B, t.C, t.D, t.E
}

func Foo(a, b, c string, d ...string) {
	println(a, b, c)
}

func Bar(t Tuple5[string, string, string, string, string]) {

	Foo(t.Deconstruct())
}
