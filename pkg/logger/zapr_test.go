package logger

import (
	"fmt"
	"testing"
)

func TestZapr(t *testing.T) {
	// var log logr.Logger

	// zapLog, err := zap.NewDevelopment()
	// if err != nil {
	// 	panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	// }

	log := NewZapr()

	log.V(3).Info("3")
	log.V(2).Info("2")

	log.V(1).Info("1")

	log.V(0).Info("0")

	log.V(-1).Info("-1")

	log.V(-2).Info("-2")
	log.V(-2).Error(fmt.Errorf("error"), "error")

	withName := log.WithName("foo").WithName("bar")

	withName.Info("withName")
}
