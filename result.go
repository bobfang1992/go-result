package result

type Result[T any] struct {
	value T
	err error
}

func (r Result[T]) Ok() bool {
	return r.err == nil
}

func (r Result[T]) ValueOr(v T) T {
	if r.Ok() {
		return r.value
	}
	return v
}

func (r Result[T]) ValueOrPanic() T {
	if r.Ok() {
		return r.value
	}
	panic(r.err)
}