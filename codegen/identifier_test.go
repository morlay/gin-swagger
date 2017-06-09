package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitToWords(t *testing.T) {
	assert.Equal(t, []string{}, SplitToWords(""))
	assert.Equal(t, []string{"lowercase"}, SplitToWords("lowercase"))
	assert.Equal(t, []string{"Class"}, SplitToWords("Class"))
	assert.Equal(t, []string{"My", "Class"}, SplitToWords("MyClass"))
	assert.Equal(t, []string{"My", "C"}, SplitToWords("MyC"))
	assert.Equal(t, []string{"HTML"}, SplitToWords("HTML"))
	assert.Equal(t, []string{"PDF", "Loader"}, SplitToWords("PDFLoader"))
	assert.Equal(t, []string{"A", "String"}, SplitToWords("AString"))
	assert.Equal(t, []string{"Simple", "XML", "Parser"}, SplitToWords("SimpleXMLParser"))
	assert.Equal(t, []string{"vim", "RPC", "Plugin"}, SplitToWords("vimRPCPlugin"))
	assert.Equal(t, []string{"GL", "11", "Version"}, SplitToWords("GL11Version"))
	assert.Equal(t, []string{"99", "Bottles"}, SplitToWords("99Bottles"))
	assert.Equal(t, []string{"May", "5"}, SplitToWords("May5"))
	assert.Equal(t, []string{"BFG", "9000"}, SplitToWords("BFG9000"))
	assert.Equal(t, []string{"Böse", "Überraschung"}, SplitToWords("BöseÜberraschung"))
	assert.Equal(t, []string{"Two", "spaces"}, SplitToWords("Two  spaces"))
	assert.Equal(t, []string{"BadUTF8\xe2\xe2\xa1"}, SplitToWords("BadUTF8\xe2\xe2\xa1"))
	assert.Equal(t, []string{"snake", "case"}, SplitToWords("snake_case"))
	assert.Equal(t, []string{"snake", "case"}, SplitToWords("snake_ case"))
}

func TestToUpperCamelCase(t *testing.T) {
	assert.Equal(t, "SnakeCase", ToUpperCamelCase("snake_case"))
	assert.Equal(t, "IDCase", ToUpperCamelCase("id_case"))
}

func TestToLowerCamelCase(t *testing.T) {
	assert.Equal(t, "snakeCase", ToLowerCamelCase("snake_case"))
	assert.Equal(t, "idCase", ToLowerCamelCase("id_case"))
}

func TestToUpperSnakeCase(t *testing.T) {
	assert.Equal(t, "SNAKE_CASE", ToUpperSnakeCase("snakeCase"))
	assert.Equal(t, "ID_CASE", ToUpperSnakeCase("idCase"))
}

func TestToLowerSnakeCase(t *testing.T) {
	assert.Equal(t, "snake_case", ToLowerSnakeCase("snakeCase"))
	assert.Equal(t, "id_case", ToLowerSnakeCase("idCase"))
	assert.Equal(t, "i7_case", ToLowerSnakeCase("i7Case"))
}
