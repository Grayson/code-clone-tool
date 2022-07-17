package stage

type Stage[T any] struct {
	Value T
	Err   error
}

func Start[T any](f func() (T, error)) Stage[T] {
	v, e := f()
	return Stage[T]{
		v,
		e,
	}
}

type mapNext[T any, U any] func(T) (U, error)

func Then[T any, U any](prev Stage[T], next mapNext[T, U]) Stage[U] {
	v, e := Finally(prev, next)
	return Stage[U]{
		v,
		e,
	}
}

func Finally[T any, U any](prev Stage[T], next mapNext[T, U]) (U, error) {
	var v U
	if prev.Err != nil {
		return v, prev.Err
	}
	return next(prev.Value)
}
