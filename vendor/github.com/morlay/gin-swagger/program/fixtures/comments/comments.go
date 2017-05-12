// package comments
package comments

// import
import (
	// fmt
	"fmt"
	"time"
)

// swagger:strfmt date-time
type Date time.Time

// Struct
type Test struct {
	// TestFieldString
	String string
	// TestFieldInt
	Int int
	// TestFieldBool
	Bool bool
	// TestFieldDate
	Date Date
}

// TypeGroup
type (
	Test2 struct {
		// TestFieldString
		String string
		// TestFieldInt
		Int int
		// TestFieldBool
		Bool bool
	}
)

// Test
func (t Test2) Recv() {

}

// SomeVar
var test = Test{
	String: "",
	Int:    1 + 1,
	Bool:   true,
}

// VarGroup
var (
	// test2
	test2 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
	// test2
	test3 = Test{
		String: "",
		Int:    1,
		Bool:   true,
	}
)

// Print
func Print(a string, b string) string {
	return a + b
}

// SomeFunc
func fn() {
	// Call
	res := Print("", "")
	if res != "" {
		// print
		fmt.Println(res)
	}
}
