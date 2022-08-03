package lib

import "strings"

type BoolString string

func NewBoolString(b bool) BoolString {
	if b {
		return "true"
	}
	return "false"
}

func (b BoolString) IsTruthy() bool {
	switch strings.ToLower(string(b)) {
	case "y":
		fallthrough
	case "yes":
		fallthrough
	case "t":
		fallthrough
	case "true":
		fallthrough
	case "1":
		return true
	}
	return false
}
