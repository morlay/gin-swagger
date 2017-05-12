package swagger_to_client

import (
	"fmt"
	"sort"

	"github.com/go-openapi/spec"
	"github.com/morlay/gin-swagger/codegen"
)

func getFieldsFromSchema(schema spec.Schema) (fields []string, deps []string) {
	var propNames = []string{}

	for name := range schema.Properties {
		propNames = append(propNames, name)
	}

	sort.Strings(propNames)

	for _, name := range propNames {
		propSchema := schema.Properties[name]
		fieldName := codegen.ToUpperCamelCase(name)

		if propSchema.Extensions["x-go-name"] != nil {
			fieldName = fmt.Sprint(propSchema.Extensions["x-go-name"])
		}

		goType, subDeps := GetTypeFromSchema(propSchema)
		deps = append(deps, subDeps...)

		var tags []string

		var jsonTag = name

		if propSchema.Enum == nil && propSchema.Type.Contains("string") && goType != "string" {
			jsonTag = codegen.JoinWithComma(jsonTag, "string")
		}

		tags = append(tags, codegen.DeclTag("json", jsonTag))

		if fmt.Sprint(schema.Default) != "<nil>" {
			tags = append(tags, codegen.DeclTag("default", fmt.Sprint(schema.Default)))
		}

		if propSchema.Extensions["x-go-validate"] != nil {
			tags = append(tags, codegen.DeclTag("validate", fmt.Sprint(propSchema.Extensions["x-go-validate"])))
		}

		fields = append(fields, codegen.DeclField(
			fieldName,
			goType,
			tags,
			propSchema.Description,
		))
	}
	return
}

func GetTypeFromSchema(schema spec.Schema) (tpe string, deps []string) {
	if schema.Ref.String() != "" {
		tpe = getRefName(schema.Ref.String())
		return
	}

	if schema.Extensions["x-go-named"] != nil {
		tpeName := fmt.Sprint(schema.Extensions["x-go-named"])
		tpe = getRefName(tpeName)
		deps = append(deps, getPackageNameFromPath(tpeName))
		return
	}

	if len(schema.AllOf) > 0 {
		var fields []string

		for _, subSchema := range schema.AllOf {
			if subSchema.Ref.String() != "" {
				gType := getRefName(subSchema.Ref.String())
				fields = append(fields, codegen.DeclField(
					"",
					gType,
					[]string{""},
					"",
				))
			}

			if subSchema.Properties != nil {
				otherFields, subDeps := getFieldsFromSchema(subSchema)
				fields = append(fields, otherFields...)
				deps = append(deps, subDeps...)
			}
		}

		tpe = codegen.DeclStruct(fields)
		return

	}

	if schema.Type.Contains("object") {
		if schema.AdditionalProperties != nil {
			goType, subDeps := GetTypeFromSchema(*schema.AdditionalProperties.Schema)
			deps = append(deps, subDeps...)
			tpe = codegen.DeclMap("string", goType)
			return
		}

		if schema.Properties != nil {
			fields, subDeps := getFieldsFromSchema(schema)
			deps = append(deps, subDeps...)

			tpe = codegen.DeclStruct(fields)
			return
		}

	}

	if schema.Type.Contains("array") {
		if schema.Items != nil {
			goType, subDeps := GetTypeFromSchema(*schema.Items.Schema)
			deps = append(deps, subDeps...)
			tpe = codegen.DeclSlice(goType)
			return
		}
	}

	schemaType := schema.Type[0]
	format := schema.Format

	switch format {
	case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		tpe = format
	case "float":
		tpe = "float32"
	case "double":
		tpe = "float64"
	default:
		switch schemaType {
		case "boolean":
			tpe = "bool"
		default:
			tpe = "string"
		}
	}

	return
}

func ToGoType(name string, schema spec.Schema) (string, []string) {
	goType, deps := GetTypeFromSchema(schema)

	return codegen.DeclType(name, goType), deps
}

func ToTypes(pkgName string, swagger spec.Swagger) string {
	p := codegen.NewPrinter().Input(codegen.DeclPackage(pkgName)).NewLine()

	var types = []string{}
	var deps = []string{}
	var definitionNames = []string{}

	for name := range swagger.Definitions {
		definitionNames = append(definitionNames, name)
	}

	sort.Strings(definitionNames)

	for _, name := range definitionNames {
		goType, subDeps := ToGoType(name, swagger.Definitions[name])
		types = append(types, goType)
		deps = append(deps, subDeps...)
	}

	p.Input(codegen.DeclImports(deps...)).NewLine()
	p.Input(codegen.JoinWithLineBreak(types...))

	return p.String()
}
