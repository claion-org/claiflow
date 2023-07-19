package macro

import (
	"fmt"
	"testing"
	"time"
)

func TestIntervalLogE(t *testing.T) {
	var acc = time.Duration(0)
	for i := 0; i < 100; i++ {

		got := IntervalSqrt(acc, time.Second, (time.Second * 3), i)
		t.Log(i, got)
		fmt.Println(i, "\t", float64(got))
		acc = got
	}
}
