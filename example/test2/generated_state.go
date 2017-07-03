package test2

import (
	"errors"
	"strings"
)

var InvalidState = errors.New("invalid state")

func (v State) String() string {
	switch v {
	case STATE_UNKNOWN:
		return ""
	case STATE__ONE:
		return "ONE"
	case STATE__TWO:
		return "TWO"
	case STATE__THREE:
		return "THREE"
	case STATE__FOUR:
		return "FOUR"
	}
	return "UNKNOWN"
}

func (v State) Label() string {
	switch v {
	case STATE_UNKNOWN:
		return ""
	case STATE__ONE:
		return "one"
	case STATE__TWO:
		return "two"
	case STATE__THREE:
		return "three"
	case STATE__FOUR:
		return "four"
	}
	return "UNKNOWN"
}

func ParseStateFromString(s string) (State, error) {
	switch s {
	case "":
		return STATE_UNKNOWN, nil
	case "ONE":
		return STATE__ONE, nil
	case "TWO":
		return STATE__TWO, nil
	case "THREE":
		return STATE__THREE, nil
	case "FOUR":
		return STATE__FOUR, nil
	}
	return STATE_UNKNOWN, InvalidState
}

func (v State) MarshalJSON() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidState
	}
	return []byte("\"" + str + "\""), nil
}

func (v *State) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(strings.ToUpper(string(data)), "\"")
	*v, err = ParseStateFromString(s)
	return
}
