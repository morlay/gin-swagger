package test2

import (
	"errors"
	"strings"
)

var InvalidBool = errors.New("invalid bool")

func (v Bool) String() string {
	switch v {
	case BOOL_UNKNOWN:
		return ""
	case BOOL__TRUE:
		return "TRUE"
	case BOOL__FALSE:
		return "FALSE"
	}
	return "UNKNOWN"
}

func ParseBoolFromString(s string) (Bool, error) {
	switch s {
	case "":
		return BOOL_UNKNOWN, nil
	case "TRUE":
		return BOOL__TRUE, nil
	case "FALSE":
		return BOOL__FALSE, nil
	}
	return BOOL_UNKNOWN, InvalidBool
}

func (v Bool) MarshalJSON() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidBool
	}
	return []byte("\"" + str + "\""), nil
}

func (v *Bool) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(strings.ToUpper(string(data)), "\"")
	*v, err = ParseBoolFromString(s)
	return
}
