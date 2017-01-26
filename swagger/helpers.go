package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
	"regexp"
	"strconv"
	"strings"
)

var (
	rxEnum     = regexp.MustCompile(`swagger:enum`)
	rxStrFmt   = regexp.MustCompile(`swagger:strfmt\s+(\S+)([\s\S]+)?$`)
	rxValidate = regexp.MustCompile(`^@([^\[]+)\[([^\]]+)\]$`)
)

func ParseEnum(str string) (string, bool) {
	if rxEnum.MatchString(str) {
		return strings.TrimSpace(strings.Replace(str, "swagger:enum", "", -1)), true
	}

	return str, false
}

func ParseStrfmt(str string) (string, string) {
	matched := rxStrFmt.FindAllStringSubmatch(str, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return str, ""
}

func exclusiveValues(str string) (string, bool) {
	var values = strings.SplitN(str, "^", -1)
	if len(values) == 2 {
		return values[1], true
	}
	return values[0], false
}

func GetCommonValidations(validateTag string) (commonValidations spec.CommonValidations) {
	var matched = rxValidate.FindAllStringSubmatch(validateTag, -1)

	if len(matched) > 0 && len(matched[0]) == 3 {
		tpe := matched[0][1]
		params := strings.Split(matched[0][2], ",")
		switch tpe {
		case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
			if len(params) > 0 {
				value, exclusive := exclusiveValues(params[0])
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					commonValidations.Minimum = &val
					commonValidations.ExclusiveMinimum = exclusive
				}
			}
			if len(params) > 1 {
				value, exclusive := exclusiveValues(params[1])
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					commonValidations.Maximum = &val
					commonValidations.ExclusiveMaximum = exclusive
				}
			}
		case "string":
			if len(params) > 0 {
				if val, err := strconv.ParseInt(params[0], 10, 64); err == nil {
					commonValidations.MinLength = &val
				}
			}
			if len(params) > 1 {
				if val, err := strconv.ParseInt(params[1], 10, 64); err == nil {
					commonValidations.MaxLength = &val
				}
			}
		}
	}

	return
}

func GetSchemaTypeFromBasicType(basicTypeName string) (string, string, bool) {
	switch basicTypeName {
	case "bool":
		return "boolean", "", true
	case "byte":
		return "integer", "uint8", true
	case "error":
		return "string", "", true
	case "float32":
		return "number", "float", true
	case "float64":
		return "number", "double", true
	case "int":
		return "integer", "int64", true
	case "int8":
		return "integer", "int8", true
	case "int16":
		return "integer", "int16", true
	case "int32":
		return "integer", "int32", true
	case "int64":
		return "integer", "int64", true
	case "rune":
		return "integer", "int32", true
	case "string":
		return "string", "", true
	case "uint":
		return "integer", "uint64", true
	case "uint16":
		return "integer", "uint16", true
	case "uint32":
		return "integer", "uint32", true
	case "uint64":
		return "integer", "uint64", true
	case "uint8":
		return "integer", "uint8", true
	case "uintptr":
		return "integer", "uint64", true
	default:
		panic(fmt.Errorf("unsupported type %q", basicTypeName))
	}
	return "", "", false
}

func BindSchemaWithCommonValidations(schema *spec.Schema, commonValidations spec.CommonValidations) {
	schema.Maximum = commonValidations.Maximum
	schema.ExclusiveMaximum = commonValidations.ExclusiveMaximum
	schema.Minimum = commonValidations.Minimum
	schema.ExclusiveMinimum = commonValidations.ExclusiveMinimum
	schema.MaxLength = commonValidations.MaxLength
	schema.MinLength = commonValidations.MinLength
	schema.Pattern = commonValidations.Pattern
	schema.MaxItems = commonValidations.MaxItems
	schema.MinItems = commonValidations.MinItems
	schema.UniqueItems = commonValidations.UniqueItems
	schema.MultipleOf = commonValidations.MultipleOf
	schema.Enum = commonValidations.Enum
}

func BindCommonValidationsWithSchema(commonValidations *spec.CommonValidations, schema spec.Schema) {
	commonValidations.Maximum = schema.Maximum
	commonValidations.ExclusiveMaximum = schema.ExclusiveMaximum
	commonValidations.Minimum = schema.Minimum
	commonValidations.ExclusiveMinimum = schema.ExclusiveMinimum
	commonValidations.MaxLength = schema.MaxLength
	commonValidations.MinLength = schema.MinLength
	commonValidations.Pattern = schema.Pattern
	commonValidations.MaxItems = schema.MaxItems
	commonValidations.MinItems = schema.MinItems
	commonValidations.UniqueItems = schema.UniqueItems
	commonValidations.MultipleOf = schema.MultipleOf
	commonValidations.Enum = schema.Enum
}

func BindSimpleSchemaWithSchema(simpleSchema *spec.SimpleSchema, schema spec.Schema) {
	simpleSchema.Type = schema.Type[0]
	simpleSchema.Format = schema.Format
	simpleSchema.Default = schema.Default
}

func BindParameterWithSchema(param *spec.Parameter, schema spec.Schema) {
	param.VendorExtensible = schema.VendorExtensible
	BindSimpleSchemaWithSchema(&param.SimpleSchema, schema)
	BindCommonValidationsWithSchema(&param.CommonValidations, schema)
}

func BindItemsWithSchema(items *spec.Items, schema spec.Schema) {
	items.VendorExtensible = schema.VendorExtensible
	BindSimpleSchemaWithSchema(&items.SimpleSchema, schema)
	BindCommonValidationsWithSchema(&items.CommonValidations, schema)
}
