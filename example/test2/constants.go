package test2

import (
	"errors"
)

// swagger:enum
type (
	State int
)

const (
	STATE_UNKNOWN = iota
	STATE__ONE    // one
	STATE__TWO    // two
	STATE__THREE  // two
)

var (
	StateKeyValueMap = map[State]string{
		STATE_UNKNOWN: "",
		STATE__ONE:    "NORMAL",
		STATE__TWO:    "FROZEN",
		STATE__THREE:  "THREE",
	}
	StateValueKeyMap = RevertToValueKey(StateKeyValueMap)
	InvalidState     = errors.New("invalid State")
)

func RevertToValueKey(kvMap map[State]string) map[string]State {
	vkMap := map[string]State{}

	for key, value := range kvMap {
		vkMap[value] = key
	}
	return vkMap
}

func (state State) MarshalJSON() ([]byte, error) {
	if value, ok := StateKeyValueMap[state]; ok {
		return []byte(`"` + value + `"`), nil
	}
	return nil, InvalidState
}

func (state *State) UnmarshalJSON(data []byte) error {
	if key, ok := StateValueKeyMap[string(data)]; ok {
		*state = key
		return nil
	}
	return InvalidState
}
