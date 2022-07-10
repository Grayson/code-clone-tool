package optional

type Optional[T any] struct {
	value    T
	HasValue bool
}

func Of[T any](value T) *Optional[T] {
	return &Optional[T]{
		value,
		true,
	}
}

func OfNonZero[T comparable](value T) *Optional[T] {
	var zeroValue T
	return &Optional[T]{
		value,
		value != zeroValue,
	}
}

func Empty[T comparable]() *Optional[T] {
	var zero T
	return OfNonZero(zero)
}

func (opt *Optional[T]) Get() (v T, ok bool) {
	ok = opt.HasValue
	if ok {
		v = opt.value
	}
	return
}
