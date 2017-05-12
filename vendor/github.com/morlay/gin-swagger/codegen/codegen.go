package codegen

import "strings"

func DeclPackage(name string) string {
	return NewPrinter().Input(JoinWithSpace("package", name)).NewLine().String()
}

func DeclType(name string, tpe string) string {
	return NewPrinter().Input(JoinWithSpace("type", name, tpe)).NewLine().String()
}

func DeclVar(name string, value string) string {
	return NewPrinter().Input(JoinWithSpace("var", name, "=", value)).NewLine().String()
}

func DeclMap(keyType string, valueType string) string {
	return "map" + WithSquareBrackets(keyType) + valueType
}

func DeclReturn(value string) string {
	return "return " + value
}

func DeclCase(value string) string {
	return "case " + value + ":"
}

func DeclSlice(itemType string) string {
	return WithSquareBrackets("") + itemType
}

func DeclTag(key string, value string) string {
	return JoinWithColon(key, WithQuotes(value))
}

func DeclField(name string, tpe string, tags []string, desc string) string {
	if name == "" {
		return tpe
	}

	var comments []string

	for _, s := range strings.Split(desc, "\n") {
		comments = append(comments, "// "+s)
	}

	return JoinWithLineBreak(
		JoinWithLineBreak(comments...),
		JoinWithSpace(
			name,
			tpe,
			WithTagQuotes(JoinWithSpace(tags...)),
		),
	)
}

func DeclStruct(fields []string) string {
	return NewPrinter().Input("struct").Space().
		Input(WithCurlyBrackets(
			JoinWithLineBreak(
				"",
				JoinWithLineBreak(fields...),
				"",
			))).
		String()
}

func DeclImports(deps ...string) string {
	if len(deps) > 0 {
		return "import " + WithRoundBrackets(
			JoinWithLineBreak(
				"",
				JoinWithLineBreak(
					MapStrings(WithQuotes, UniqueStrings(deps))...,
				),
				"",
			),
		)
	}
	return ""
}
