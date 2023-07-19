package macro_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func Test_Errorf(t *testing.T) {
	err := fmt.Errorf("%w: test error message", sql.ErrNoRows)

	t.Log(err)
}

func Test_Unwrap(t *testing.T) {
	{
		err := fmt.Errorf("%w: test error message", sql.ErrNoRows)

		orig := errors.Unwrap(err)

		ok := sql.ErrNoRows == orig

		if ok != true {
			t.Errorf("expected=%v", true)
		}
	}
	{
		err := fmt.Errorf("%w: test error message", sql.ErrNoRows)

		orig := errors.Unwrap(err)

		ok := sql.ErrTxDone == orig
		if ok != false {
			t.Errorf("expected=%v", false)
		}
	}
	{
		err1 := fmt.Errorf("test error message")
		err2 := fmt.Errorf("test error message")

		if (err1 == err2) != false {
			t.Errorf("expected=%v", false)
		}

		if (err1.Error() == err2.Error()) != true {
			t.Errorf("expected=%v", true)
		}
	}

}

func Test_Is(t *testing.T) {
	{
		err := fmt.Errorf("%w: test error message", sql.ErrNoRows)

		ok := errors.Is(err, sql.ErrNoRows)

		if ok != true {
			t.Errorf("expected=%v", true)
		}
	}
}
