package lib

import "testing"

func TestBoolString_IsTruthy(t *testing.T) {
	tests := []struct {
		name string
		b    BoolString
		want bool
	}{
		{
			"`y` is true",
			"y",
			true,
		},
		{
			"`yes` is true",
			"yes",
			true,
		},
		{
			"`t` is true",
			"t",
			true,
		},
		{
			"`true` is true",
			"true",
			true,
		},
		{
			"`1` is true",
			"1",
			true,
		},
		{
			"`0` is false",
			"0",
			false,
		},
		{
			"`n` is false",
			"n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsTruthy(); got != tt.want {
				t.Errorf("BoolString.IsTruthy() = %v, want %v", got, tt.want)
			}
		})
	}
}
