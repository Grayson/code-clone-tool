package either

type inner interface{ hasValue() bool }
type filled[T any] struct{ Value T }
type empty struct{}

func (f filled[T]) hasValue() bool { return true }
func (e empty) hasValue() bool     { return false }

type Either[T any, U any] struct {
	l inner
	r inner
}

func Of[T any, U any](value any) *Either[T, U] {
	switch v := value.(type) {
	case T:
		return &Either[T, U]{
			filled[T]{v},
			empty{},
		}
	case U:
		return &Either[T, U]{
			empty{},
			filled[U]{v},
		}
	}
	return nil
}

func (of *Either[T, U]) GetLeft() (out T, ok bool) {
	ok = false
	if of.l.hasValue() {
		out = of.l.(filled[T]).Value
		ok = true
	}
	return
}

func (of *Either[T, U]) GetRight() (out U, ok bool) {
	ok = false
	if of.r.hasValue() {
		out = of.r.(filled[U]).Value
		ok = true
	}
	return
}
