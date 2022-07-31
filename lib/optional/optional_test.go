package optional_test

import (
	"testing"

	"github.com/grayson/code-clone-tool/lib/optional"
)

func TestEmpty(t *testing.T) {
	o := optional.Empty[*int]()
	if o.HasValue {
		t.Error("Optional unexpected has value")
	}
	v, ok := o.Get()
	if v != nil {
		t.Error("Expected nil value in optional")
	}

	if ok {
		t.Error("Expected `ok` to be false")
	}
}

func TestFilled(t *testing.T) {
	i := 42
	o := optional.Of(&i)
	if !o.HasValue {
		t.Error("Expected optional to have value")
	}
	v, ok := o.Get()
	if v == nil {
		t.Error("Expected optional to have valid pointer")
	}

	deref := *v
	if deref != 42 {
		t.Errorf("Unexpected value in optional: %d (expected %d)", deref, 42)
	}

	if !ok {
		t.Error("Expected `ok` to be true")
	}
}

func TestValueType(t *testing.T) {
	o := optional.Of(42)
	if !o.HasValue {
		t.Error("Expected optional to have value")
	}
	v, ok := o.Get()
	if v != 42 {
		t.Errorf("Unexpected value in optional: %d (expected %d)", v, 42)
	}

	if !ok {
		t.Error("Expected `ok` to be true")
	}
}

func TestValueTypeInNillableApi(t *testing.T) {
	o := optional.OfNonZero(42)
	if !o.HasValue {
		t.Error("Expected optional to have value")
	}
	v, ok := o.Get()
	if v != 42 {
		t.Errorf("Unexpected value in optional: %d (expected %d)", v, 42)
	}

	if !ok {
		t.Error("Expected `ok` to be true")
	}
}

func TestZeroValueOfValueTypeInNonZeroApi(t *testing.T) {
	o := optional.OfNonZero(0)
	if o.HasValue {
		t.Error("Expected optional to not have value")
	}
	v, ok := o.Get()
	if v != 0 {
		t.Errorf("Unexpected value in optional: %d (expected %d)", v, 0)
	}

	if ok {
		t.Error("Expected `ok` to be false")
	}
}
