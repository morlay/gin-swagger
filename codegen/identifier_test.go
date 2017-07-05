package codegen_test

import (
	"testing"

	"github.com/morlay/gin-swagger/codegen"
	"github.com/stretchr/testify/assert"
)

func TestSplitToWords(tt *testing.T) {
	t := assert.New(tt)

	t.Equal([]string{}, codegen.SplitToWords(""))
	t.Equal([]string{"lowercase"}, codegen.SplitToWords("lowercase"))
	t.Equal([]string{"Class"}, codegen.SplitToWords("Class"))
	t.Equal([]string{"My", "Class"}, codegen.SplitToWords("MyClass"))
	t.Equal([]string{"My", "C"}, codegen.SplitToWords("MyC"))
	t.Equal([]string{"HTML"}, codegen.SplitToWords("HTML"))
	t.Equal([]string{"PDF", "Loader"}, codegen.SplitToWords("PDFLoader"))
	t.Equal([]string{"A", "String"}, codegen.SplitToWords("AString"))
	t.Equal([]string{"Simple", "XML", "Parser"}, codegen.SplitToWords("SimpleXMLParser"))
	t.Equal([]string{"vim", "RPC", "Plugin"}, codegen.SplitToWords("vimRPCPlugin"))
	t.Equal([]string{"GL", "11", "Version"}, codegen.SplitToWords("GL11Version"))
	t.Equal([]string{"99", "Bottles"}, codegen.SplitToWords("99Bottles"))
	t.Equal([]string{"May", "5"}, codegen.SplitToWords("May5"))
	t.Equal([]string{"BFG", "9000"}, codegen.SplitToWords("BFG9000"))
	t.Equal([]string{"Böse", "Überraschung"}, codegen.SplitToWords("BöseÜberraschung"))
	t.Equal([]string{"Two", "spaces"}, codegen.SplitToWords("Two  spaces"))
	t.Equal([]string{"BadUTF8\xe2\xe2\xa1"}, codegen.SplitToWords("BadUTF8\xe2\xe2\xa1"))
	t.Equal([]string{"snake", "case"}, codegen.SplitToWords("snake_case"))
	t.Equal([]string{"snake", "case"}, codegen.SplitToWords("snake_ case"))
}

func TestToUpperCamelCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("SnakeCase", codegen.ToUpperCamelCase("snake_case"))
	t.Equal("IDCase", codegen.ToUpperCamelCase("id_case"))
}

func TestToLowerCamelCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("snakeCase", codegen.ToLowerCamelCase("snake_case"))
	t.Equal("idCase", codegen.ToLowerCamelCase("id_case"))
}

func TestToUpperSnakeCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("SNAKE_CASE", codegen.ToUpperSnakeCase("snakeCase"))
	t.Equal("ID_CASE", codegen.ToUpperSnakeCase("idCase"))
}

func TestToLowerSnakeCase(tt *testing.T) {
	t := assert.New(tt)

	t.Equal("snake_case", codegen.ToLowerSnakeCase("snakeCase"))
	t.Equal("id_case", codegen.ToLowerSnakeCase("idCase"))
	t.Equal("i7_case", codegen.ToLowerSnakeCase("i7Case"))
}
