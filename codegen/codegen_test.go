package codegen

import (
	"fmt"
	"testing"
)

func TestPrinter(t *testing.T) {
	fmt.Println(DeclPackage("some_package"))
	fmt.Println(DeclType("Test", "int"))
}
