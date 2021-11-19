package result

type Result[T any] struct {
	value T
	err   error
}

// Value related operations

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

// Monadic operations

// Map and then are actually the same function, different languages use different names for the same thing
// Since this is an expeiremental implemenation, I will provide both names so we can deceide which one to use
// There is some limitation to Golang's Generics at this moment, so we cannot havea a fully fledged Map functon
// like Haskell.
// This function has the following type signature in Haskell: (>>=) :: Maybe a -> (a -> Maybe b) -> Maybe b
// But in Golang, we currently do not support having another idependent type variable in method signature, so we are
// limited to use T here, but in reality we should have used func(T) Result[U] for type of f.
func (r Result[T]) Map(f func(T) Result[T]) Result[T] {
	if r.Ok() {
		return f(r.value)
	}
	return r
}

func (r Result[T]) Then(f func(T) Result[T]) Result[T] {
	return r.Map(f)
}
