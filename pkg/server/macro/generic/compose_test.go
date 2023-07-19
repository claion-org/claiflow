package generic

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCompose(t *testing.T) {

	AtoI := Compose(Placeholder[string](), strconv.Atoi)
	ItoA := Compose(AtoI, Wrap(strconv.Itoa))
	StrLen := Compose(Placeholder[string](), func(s string) (int, error) { return len(s), nil })

	a, err := ItoA("10")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a)

	i, err := AtoI("20")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(i)

	_, err = AtoI("20A")
	if err != nil {
		t.Error(err)
	}

	l, err := StrLen("foo")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(l)
}
