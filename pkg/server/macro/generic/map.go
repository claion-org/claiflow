package generic

func MapToArray[A comparable, B any](mAB map[A]B) []B {
	bb := make([]B, len(mAB))
	for _, v := range mAB {
		bb = append(bb, v)
	}

	return bb
}

func ArrayToMap[A comparable, B any](fn func(B) A, bb []B) map[A]B {
	mAB := make(map[A]B)
	for _, b := range bb {
		a := fn(b)
		mAB[a] = b
	}

	return mAB
}
