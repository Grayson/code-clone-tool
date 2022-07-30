package either_test

import (
	"testing"

	"github.com/Grayson/code-clone-tool/lib/either"
)

func TestBasicOfCreation(t *testing.T) {
	of := either.Of[int, bool](42)
	if of == nil {
		t.Error("of was nil")
	}
}

func TestLeftValueGetter(t *testing.T) {
	of := either.Of[int, bool](42)
	v, ok := of.GetLeft()
	if v != 42 {
		t.Errorf("Received wrong input: %v (expected %v)", v, 42)
	}
	if !ok {
		t.Error("Expected ok of left to be true")
	}

	if _, ok := of.GetRight(); ok {
		t.Error("Expected ok of right to be false")
	}
}

func TestRightValueGetter(t *testing.T) {
	of := either.Of[int, bool](true)
	v, ok := of.GetRight()
	if !v {
		t.Errorf("Received wrong input: %v (expected %v)", v, true)
	}
	if !ok {
		t.Error("Expected ok of right to be true")
	}

	if _, ok := of.GetLeft(); ok {
		t.Error("Expected ok of left to be false")
	}
}
