package swagger_to_client

import (
	"encoding/json"

	"github.com/go-openapi/spec"
	"github.com/morlay/gin-swagger/codegen"
)

func NewClientGenerator(name string, baseClient string) *ClientGenerator {
	return &ClientGenerator{
		Name:       name,
		BaseClient: baseClient,
	}
}

type ClientGenerator struct {
	BaseClient string
	Name       string
	Swagger    spec.Swagger
}

func (c *ClientGenerator) LoadSwagger(swagger spec.Swagger) {
	c.Swagger = swagger
}

func (c *ClientGenerator) LoadSwaggerFromFile(path string) {
	c.Swagger = spec.Swagger{}

	swaggerString := codegen.OpenFile(path)
	err := json.Unmarshal([]byte(swaggerString), &c.Swagger)
	if err != nil {
		panic(err)
	}
}

func (c *ClientGenerator) Output() {
	pkgName := codegen.ToLowerSnakeCase("Client-" + c.Name)
	codegen.GenerateGoFile(codegen.JoinWithSlash(pkgName, "types.go"), ToTypes(pkgName, c.Swagger))
	codegen.GenerateGoFile(codegen.JoinWithSlash(pkgName, "client.go"), ToClient(c.BaseClient, pkgName, c.Swagger))
}
