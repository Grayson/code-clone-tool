package stage

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestStart(t *testing.T) {
	type args struct {
		f func() (string, error)
	}
	type want struct {
		s string
		e error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Start with value",
			args{func() (string, error) { return "test", nil }},
			want{"test", nil},
		},
		{
			"Start with error",
			args{func() (string, error) { return "", fmt.Errorf("error") }},
			want{"", fmt.Errorf("error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Start(tt.args.f)
			correctValue := tt.want.s == got.Value
			correctErr := (got.Err == nil && tt.want.e == nil) || (got.Err != nil && tt.want.e != nil)
			if !correctValue || !correctErr {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThen(t *testing.T) {
	prev := Start(func() (string, error) {
		return "test", nil
	})
	errprev := Start(func() (string, error) {
		return "", fmt.Errorf("err")
	})
	testerr := fmt.Errorf("err")

	type args struct {
		prev Stage[string]
		next mapNext[string, string]
	}
	type want struct {
		s string
		e error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Map valid previous to next",
			args{prev, func(s string) (string, error) { return s + "n", nil }},
			want{"testn", nil},
		},
		{
			"Persist error from prev",
			args{errprev, func(s string) (string, error) { return "", nil }},
			want{"", errprev.Err},
		},
		{
			"Throw new error in next",
			args{prev, func(s string) (string, error) { return "", testerr }},
			want{"", testerr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Then(tt.args.prev, tt.args.next)
			correctValue := tt.want.s == got.Value
			correctErr := (got.Err == nil && tt.want.e == nil) || (got.Err != nil && tt.want.e != nil)
			if !correctValue || !correctErr {
				t.Errorf("Then() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFinally(t *testing.T) {
	validStart := Start(func() (string, error) { return "test", nil })
	validNext := Then(
		validStart,
		func(s string) (string, error) { return s + "n", nil },
	)

	testerr := fmt.Errorf("err")
	invalidStart := Start(func() (string, error) { return "", testerr })
	invalidFromStartNext := Then(
		invalidStart,
		func(s string) (string, error) { return "", nil },
	)

	invalidIntermediate := Then(
		validStart,
		func(s string) (string, error) { return "", testerr },
	)

	type args struct {
		prev Stage[string]
		next mapNext[string, string]
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid happy path",
			args{validNext, func(s string) (string, error) { return s + ":)", nil }},
			"testn:)",
			false,
		},
		{
			"Invalid happy path",
			args{invalidFromStartNext, func(s string) (string, error) { return s + ":)", nil }},
			"",
			true,
		},
		{
			"Invalid intermediate path",
			args{invalidIntermediate, func(s string) (string, error) { return s + ":)", nil }},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Finally(tt.args.prev, tt.args.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("Finally() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Finally() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterate(t *testing.T) {
	validStart := Start(func() ([]string, error) {
		return []string{"1", "2", "42"}, nil
	})

	emptyStart := Start(func() ([]string, error) {
		return []string{}, nil
	})

	throwsErrorStart := Start(func() ([]string, error) {
		return []string{"1", "err", "42"}, nil
	})

	errorsFromTheStart := Start(func() ([]string, error) {
		return []string{"1", "2", "3"}, errors.New("expected")
	})
	mapper := func(s string) (int, error) {
		if s == "err" {
			return 0, errors.New("expected")
		}
		return strconv.Atoi(s)
	}

	type args struct {
		prev Stage[[]string]
		next mapNext[string, int]
	}
	tests := []struct {
		name     string
		args     args
		want     []int
		hasError bool
	}{
		{
			"Test valid iteration",
			args{validStart, mapper},
			[]int{1, 2, 42},
			false,
		},
		{
			"Test empty iteration",
			args{emptyStart, mapper},
			[]int{},
			false,
		},
		{
			"Test error thrown",
			args{throwsErrorStart, mapper},
			[]int{},
			true,
		},
		{
			"Test bailing on previous error",
			args{errorsFromTheStart, mapper},
			[]int{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Iterate(tt.args.prev, tt.args.next)
			expectedAndReceivedError := tt.hasError == (got.Err != nil)
			if !(expectedAndReceivedError || reflect.DeepEqual(got.Value, tt.want)) {
				t.Errorf("Iterate() = %v, want %v", got, tt.want)
			}
		})
	}
}
