---
title: A simple result type in Golang 1.18
date: 2021-11-19T05:15:31Z
description: Some simple test for the new blog system
---

With the introduction of generics in Go 1.18. One of the pain point I had with Golang has been finanly resolved.

In my day job, there are a lot of times I need to implement a "map-reduce" style algorithm using goroutines. Some thing like this:

```go
func processing() {
    works := []Work{
        {
            Name: "John",
            Age:  30,
        },
        {
            Name: "Jane",
            Age:  25,
        },
        ...
    }

    var wg sync.WaitGroup

    result := make(chan ProcessedWork, len(works))
    for _, work := range works {
        wg.Add(1)
        go func(work Work) {
            defer wg.Done()
            // do something
            newData, err := doSomething(work)
            result <- newData
        }(work)
    }

    wg.Wait()

    for r := range result {
        // combine results in some way
    }
    return ...
}
```

This may looks like all good, but there is a problem. What if one of the sub goroutine failed, what if this line actually return an error?

```go
newData, err := doSomething(work)
```

So usually you have two options, one is to simply introduce a seperate error channel, and the next is to introduce a new Result type locally and change the `result` channel here to accept the `Result` type. In this example we can do this:

```go
type Result struct {
    Data ProcessedWork
    Err  error
}
result := make(chan Result, len(works))
```

I usually opt for the second solution, but having to define this type every time is a pain. There were some suggestions to use `interface{}` to represent the Data and do type assertion when using the data, but generally I am not a big fan of using `interface{}` in Go.

Luckily we got generics in Go, so we out `Result` type can defined using this feature.

```go
type Result[T any] struct {
	value T
	err   error
}
```

Having this type is much better than using ad-hoc `Result` types for each `processing` function in the code base.

A number of useful methods can be added to the `Result` type, for example

```go
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
```

There are some obvious things I want to point out though.

1. Golang does not do sum type yet, there are proposals to add this, but I think it will come much later. Currenlty the best way to emulate a sum type in Golang is to simply use a struct and remeber which field you have used, like what we did here with the Result type. Other languages call this discrimated union if you include a special field which denote the field being used.
2. `Result` type, no matter what method you use to implement it, is not a new concept. C++ has `std::optional` since C++17. `Rust`, the golden boy of this era, has `Result<T, E>`. Haskell, where I first learned the conecpt of using types to represent the outcome of an operation that may or may not success, has `Maybe`.

Further to point 2, another thing to notice is that `Result` is actually a monad, `C++23` recently added monadic operation for `std::optional`, and Hasekll has the following functions in its stdlib since long before this conecpt is popular.

```haskell
return :: a -> Maybe a
return x  = Just x

(>>=) :: Maybe a -> (a -> Maybe b) -> Maybe b
(>>=) m g = case m of
                Nothing -> Nothing
                Just x  -> g x
```

The nice thing about these two functions, especially for the second function `(>>=)`, is that it allows you to easily combine/chain multiple operations that may or may not yield a result together without the need to keep using `if` to check if the result of last operation is `Ok` or not. I am not going to enumerate an example here, but if you are curious you can look at the Hasekll example [here](https://en.wikibooks.org/wiki/Haskell/Understanding_monads/Maybe).

But Golang is a bit lacking here, at this moment, if everything goes according to the plan, then we wont be able to have another idependent type variable in a generic type's method. So we cannot have something like this:

```go
func (r Result[T]) Then(f func(T) Result[S]) Result[S] { // <-- S is not allowed, we can only use T
	if r.Ok() {
		return f(r.value)
	}
	return r
}
```

IMO, this limtation quite serverly limited the usefulness of the result type.

Another thing I wish we had is C++ style "partial specialization" in Go's generics. For now the constraints for `Result` is `any`, but I do want to provide a function like this for user:

```go
func (r Result[T]) Eq(v T) bool {
    if r.Ok() {
        return r.value == v
    }
    return false
}
```

But since `T` is not `comparable` here, the `==` operatior will not work. If Go can provide a way to refine the constraints for some methods of a generic type, then it would be a nice feature. E.g.

```go
// here we refined the T type from any to comparable by providing a more precise constraints in the method receiver type
// now only Result that holds a T that are in the constraint comparable will have this method enabled.
func (r Result[T comparable]) Eq(v T) bool {

    if r.Ok() {
        return r.value == v
    }
    if r.Ok() {
        return r.value == v
    }
    return false
}
```

Sample code can be found [here](https://github.com/bobfang1992/go-result)
