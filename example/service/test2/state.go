package test2

// swagger:enum
type State int

const (
	STATE_UNKNOWN State = iota
	STATE__ONE          // one
	STATE__TWO          // two
	STATE__THREE        // three

	STATE__FOUR State = iota + 100 // four
)

// swagger:enum
type Bool int

const (
	BOOL_UNKNOWN Bool = iota
	BOOL__TRUE        // true
	BOOL__FALSE       // false
)
