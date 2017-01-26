package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
	"strings"
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

func (swagger *Swagger) AddDefinition(name string, schema spec.Schema) spec.Schema {
	if _, ok := swagger.Definitions[name]; !ok {
		fmt.Printf("added defination %s\n", name)
		swagger.Definitions[name] = schema
	} else {
		fmt.Errorf("duplicated definition %s", name)
	}
	return *spec.RefProperty("#/definitions/" + name)
}
