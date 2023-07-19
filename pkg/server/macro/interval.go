package macro

import (
	"math"
	"time"
)

// func IntervalSqrt(acc time.Duration, d time.Duration, max time.Duration, n int) time.Duration {
// 	n += 1

// 	return time.Duration(float64(acc) + (float64(d) / float64(n*n)))
// }

func IntervalSqrt(acc time.Duration, d time.Duration, max time.Duration, n int) time.Duration {
	acc_ := math.Max(
		math.Mod(
			float64(acc)+float64(d)/math.Sqrt(float64(n+1)),
			float64(max)),
		float64(d))
	return time.Duration(acc_)
}
