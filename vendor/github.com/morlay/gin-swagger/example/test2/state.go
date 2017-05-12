package test2

// swagger:enum
//go:generate gin-swagger enum State
type State int

const (
	STATE_UNKNOWN State = iota
	STATE__ONE          // one
	STATE__TWO          // two
	STATE__THREE        // three

	STATE__FOUR State = iota + 100 // four
)

// swagger:enum
//go:generate gin-swagger enum Bool
type Bool int

const (
	BOOL_UNKNOWN Bool = iota
	BOOL__TRUE        // true
	BOOL__FALSE       // false
)
