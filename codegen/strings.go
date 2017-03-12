package codegen

import (
	"bytes"
	"strings"
	"text/template"
)

func JoinWithSpace(a ...string) string {
	return strings.Join(a, " ")
}

func JoinWithColon(a ...string) string {
	return strings.Join(a, ":")
}

func JoinWithSlash(a ...string) string {
	return strings.Join(a, "/")
}

func JoinWithLineBreak(a ...string) string {
	return strings.Join(a, "\n")
}

func JoinWithComma(a ...string) string {
	return strings.Join(a, ",")
}

func WithQuotes(s string) string {
	return "\"" + s + "\""
}

func WithTagQuotes(s string) string {
	return "`" + s + "`"
}

func WithRoundBrackets(s string) string {
	return "(" + s + ")"
}

func WithSquareBrackets(s string) string {
	return "[" + s + "]"
}

func WithCurlyBrackets(s string) string {
	return "{" + s + "}"
}

func MapStrings(mapper func(s string) string, strs []string) []string {
	newStrings := []string{}

	for _, str := range strs {
		newStrings = append(newStrings, mapper(str))
	}

	return newStrings
}

func UniqueStrings(strs []string) []string {
	stringMap := map[string]bool{}

	for _, str := range strs {
		stringMap[str] = true
	}
	newStrings := []string{}

	for str := range stringMap {
		newStrings = append(newStrings, str)
	}

	return newStrings
}

func TemplateRender(s string) func(data interface{}) string {
	tmpl, err := template.New(s).Parse(s)

	if err != nil {
		panic(err)
	}

	return func(data interface{}) string {
		var value bytes.Buffer

		err = tmpl.Execute(&value, data)

		return value.String()
	}
}
