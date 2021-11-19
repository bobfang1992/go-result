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

func TestResultIntPointer(t *testing.T) {
	i := 1
	r := Result[*int]{value: &i}

	if !r.Ok() {
		t.Error("Result should be Ok")
	}

	if r.ValueOr(nil) != &i {
		t.Error("Result value should be equal to i")
	}

	if r.ValueOrPanic() != &i {
		t.Error("Result value should be equal to i")
	}
}

func TestResultStruct(t *testing.T) {
	type S struct {
		A int
		B string
	}
	s := S{A: 1, B: "b"}
	r := Result[S]{value: s}

	if !r.Ok() {
		t.Error("Result should be Ok")
	}

	if r.ValueOr(S{}).A != s.A {
		t.Error("Result value should be equal to s")
	}

	if r.ValueOrPanic().A != s.A {
		t.Error("Result value should be equal to s")
	}

	err := errors.New("error")
	rerr := Result[S]{err: err}

	if rerr.Ok() {
		t.Error("Result should not be Ok")
	}

	if rerr.ValueOr(S{}).A != 0 {
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

func TestResultStructPointer(t *testing.T) {
	type S struct {
		A int
		B string
	}
	s := S{A: 1, B: "b"}
	r := Result[*S]{value: &s}

	if !r.Ok() {
		t.Error("Result should be Ok")
	}

	if r.ValueOr(nil) != &s {
		t.Error("Result value should be equal to s")
	}

	if r.ValueOrPanic() != &s {
		t.Error("Result value should be equal to s")
	}
}

func TestResultMonadic(t *testing.T) {
	i := 1
	r := Result[int]{value: i}

	r2 := r.Then(func(v int) Result[int] {
		return Result[int]{value: v + 1}
	})

	if !r2.Ok() {
		t.Error("Result should be Ok")
	}

	if r2.ValueOr(0) != i+1 {
		t.Error("Result value should be equal to i+1")
	}

}
