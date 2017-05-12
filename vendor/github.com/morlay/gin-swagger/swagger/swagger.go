package swagger

import (
	"github.com/go-openapi/spec"
	"strings"
	"fmt"
	"regexp"
	"strconv"
	"reflect"
)

func NewSwagger() *Swagger {
	swagger := new(spec.Swagger)
	swagger.Swagger = "2.0"

	if swagger.Paths == nil {
		swagger.Paths = new(spec.Paths)
	}

	if swagger.Definitions == nil {
		swagger.Definitions = make(map[string]spec.Schema)
	}

	if swagger.Responses == nil {
		swagger.Responses = make(map[string]spec.Response)
	}

	sg := new(Swagger)
	sg.Swagger = swagger

	return sg
}

type Swagger struct {
	*spec.Swagger
}

func (swagger *Swagger) AddOperation(method string, path string, op *spec.Operation) {
	if swagger.Paths.Paths == nil {
		swagger.Paths.Paths = make(map[string]spec.PathItem)
	}

	paths := swagger.Paths
	pathObj := paths.Paths[path]

	switch strings.ToUpper(method) {
	case "GET":
		pathObj.Get = op
	case "POST":
		pathObj.Post = op
	case "PUT":
		pathObj.Put = op
	case "PATCH":
		pathObj.Patch = op
	case "HEAD":
		pathObj.Head = op
	case "DELETE":
		pathObj.Delete = op
	case "OPTIONS":
		pathObj.Options = op
	}

	paths.Paths[path] = pathObj
}

func (swagger *Swagger) AddDefinition(name string, schema spec.Schema) (*spec.Schema, bool) {
	s := spec.RefProperty("#/definitions/" + name)
	if _, ok := swagger.Definitions[name]; !ok {
		swagger.Definitions[name] = schema
		return s, true
	}
	return s, false
}

var (
	rxEnum = regexp.MustCompile(`swagger:enum`)
	rxStrFmt = regexp.MustCompile(`swagger:strfmt\s+(\S+)([\s\S]+)?$`)
	rxValidate = regexp.MustCompile(`^@([^\[\(\{]+)([\[\(\{])([^\}^\]^\)]+)([\}\]\)])$`)
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

	// "@int[1,2]", "int", "[", "1,2", "]"
	if len(matched) > 0 && len(matched[0]) == 5 {
		tpe := matched[0][1]
		startBracket := matched[0][2]
		endBracket := matched[0][4]
		values := strings.Split(matched[0][3], ",")

		if startBracket != "{" && endBracket != "}" {
			switch tpe {
			case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
				if len(values) > 0 {
					if val, err := strconv.ParseFloat(values[0], 64); err == nil {
						commonValidations.Minimum = &val
						commonValidations.ExclusiveMinimum = (startBracket == "(")
					}
				}
				if len(values) > 1 {
					if val, err := strconv.ParseFloat(values[1], 64); err == nil {
						commonValidations.Maximum = &val
						commonValidations.ExclusiveMaximum = (endBracket == ")")
					}
				}
			case "string":
				if len(values) > 0 {
					if val, err := strconv.ParseInt(values[0], 10, 64); err == nil {
						commonValidations.MinLength = &val
					}
				}
				if len(values) > 1 {
					if val, err := strconv.ParseInt(values[1], 10, 64); err == nil {
						commonValidations.MaxLength = &val
					}
				}
			}
		} else {
			var enums = []interface{}{}

			for _, value := range values {
				if tpe != "string" {
					if val, err := strconv.ParseInt(value, 10, 64); err == nil {
						enums = append(enums, val)
					}
				} else {
					enums = append(enums, value)
				}
			}

			commonValidations.Enum = enums
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

func EnumContainsValue(enum []interface{}, value interface{}) bool {
	var isContains = false

	for _, enumValue := range enum {
		if enumValue == value {
			isContains = true
		}
	}

	return isContains
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

	// for partial enum
	if len(schema.Enum) != 0 && len(commonValidations.Enum) != 0 {
		var enums []interface{}

		for _, enumValueOrIndex := range commonValidations.Enum {
			switch reflect.TypeOf(enumValueOrIndex).Name() {
			case "string":
				if EnumContainsValue(schema.Enum, enumValueOrIndex) {
					enums = append(enums, enumValueOrIndex)
				} else if enumValueOrIndex != "" {
					panic(fmt.Errorf("%s is not value of %s", enumValueOrIndex, schema.Enum))
				}
			default:
				if idx, ok := enumValueOrIndex.(int); ok {
					if schema.Enum[idx] != nil {
						enums = append(enums, schema.Enum[idx])
					} else if idx != 0 {
						panic(fmt.Errorf("%s is out-range of  %s", enumValueOrIndex, schema.Enum))
					}
				}

			}
		}

		schema.Enum = enums
	}
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