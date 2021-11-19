package result

import (
	"errors"
	"testing"
)

func TestResultInt(t *testing.T) {
	i := 1
	r := Result[int]{value: i}

	if !r.Ok() {
		t.Error("Result should be Ok")
	}

	if r.ValueOr(0) != i {
		t.Error("Result value should be equal to i")
	}

	if r.ValueOrPanic() != i {
		t.Error("Result value should be equal to i")
	}

	err := errors.New("error")
	rerr := Result[int]{err: err}

	if rerr.Ok() {
		t.Error("Result should not be Ok")
	}
	
	if rerr.ValueOr(0) != 0 {
		t.Error("Result value should be equal to 0")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Result should panic")
		} else if !errors.Is(r.(error), err) {
			t.Error("Result should panic with the error")
		}
	}()
	rerr.ValueOrPanic()

}
