package swagger_to_client

import (
	"fmt"
	"github.com/go-openapi/spec"
	"gopkg.in/gin-gonic/gin.v1"
	"sort"
	"strings"

	"github.com/morlay/gin-swagger/codegen"
)

type ClientInfo struct {
	BaseClient string
	PkgName    string
	Name       string
	Operations []OperationInfo
}

type OperationByID []OperationInfo

func (a OperationByID) Len() int {
	return len(a)
}

func (a OperationByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a OperationByID) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}

func (c *ClientInfo) RenderDecl() string {
	return codegen.TemplateRender(`
	func New{{ .Name }}(baseURL string, timeout time.Duration) *{{ .Name }} {
		return &{{ .Name }}{
			Client: client.Client{
				BaseURL: baseURL,
				Timeout: timeout / time.Millisecond,
			},
		}
	}

	type {{ .Name }} struct {
		client.Client
	}
	`)(c)
}

func getDataResponse(operation *spec.Operation, responses map[string]spec.Response) (*spec.Schema, string) {
	desc := ""

	if operation.Responses == nil || operation.Responses.StatusCodeResponses == nil {
		return nil, desc
	}

	var hasStringProduce = false

	for _, produce := range operation.Produces {
		hasStringProduce = (produce == gin.MIMEHTML)
	}

	var schema *spec.Schema

	for code, r := range operation.Responses.StatusCodeResponses {
		desc += r.Description + "\n"
		if code >= 200 && code < 300 {
			schema = &spec.Schema{}

			if hasStringProduce {
				schema.Typed("string", "")
			} else {
				schema.Typed("null", "")
			}

			if r.Ref.String() != "" && responses != nil {
				if responses[getRefName(r.Ref.String())].Schema != nil {
					schema = responses[getRefName(r.Ref.String())].Schema
				}
			}

			if r.Schema != nil {
				schema = r.Schema
			}
		}
	}
	return schema, desc
}

func (c *ClientInfo) AddOperation(method, path string, operation *spec.Operation, responses map[string]spec.Response) *ClientInfo {
	if c.Operations == nil {
		c.Operations = []OperationInfo{}
	}

	if respSchema, desc := getDataResponse(operation, responses); respSchema != nil {
		o := OperationInfo{
			ServiceName: c.Name,
			ID:          operation.ID,
			Method:      strings.ToUpper(method),
			Path:        convertSwaggerPathToGinPath(path),
			Parameters:  operation.Parameters,
			RespBody:    *respSchema,
			Description: desc,
		}

		c.Operations = append(c.Operations, o)
	}

	return c
}

func (c *ClientInfo) RenderOperations() (string, []string) {
	prefix := codegen.TemplateRender("func (c {{ .Name }}) ")(c)

	methods := []string{}
	deps := []string{}

	sort.Sort(OperationByID(c.Operations))

	for _, o := range c.Operations {
		descList := []string{}

		for _, desc := range strings.Split(o.Description, "\n") {
			if desc != "" {
				descList = append(descList, "// "+desc)
			}
		}

		sort.Strings(descList)

		methods = append(methods,
			codegen.JoinWithLineBreak(
				o.RenderReqDecl(),
				o.RenderRespDecl(),
				strings.Join(descList, "\n"),
				prefix+o.RenderOperationMethod(),
			))

		deps = append(deps, o.GetDeps()...)
	}

	return codegen.JoinWithLineBreak(methods...), deps
}

func (c *ClientInfo) Render() string {

	deps := []string{
		c.BaseClient,
		"time",
	}

	ops, subDeps := c.RenderOperations()

	deps = append(deps, subDeps...)

	return codegen.JoinWithLineBreak(
		codegen.DeclPackage(c.PkgName),
		codegen.DeclImports(deps...),
		c.RenderDecl(),
		ops,
	)
}

type OperationInfo struct {
	ServiceName string
	ID          string
	Method      string
	Path        string
	Parameters  []spec.Parameter
	RespBody    spec.Schema
	Description string
	deps        []string
}

func (op *OperationInfo) GetDeps() []string {
	return op.deps
}

func (op *OperationInfo) addDeps(deps ...string) *OperationInfo {
	if op.deps == nil {
		op.deps = []string{}
	}

	op.deps = append(op.deps, deps...)
	return op
}

func (op *OperationInfo) RenderReqDecl() string {
	if len(op.Parameters) > 0 {
		var fields []string

		for _, parameter := range op.Parameters {
			fieldName := codegen.ToUpperCamelCase(parameter.Name)

			if parameter.Extensions["x-go-name"] != nil {
				fieldName = parameter.Extensions["x-go-name"].(string)
			}

			var goType string
			var subDeps []string
			var inTag = parameter.In
			var jsonTag = parameter.Name

			if parameter.Schema != nil {
				inTag = "body"
				jsonTag = "body"
				goType, subDeps = GetTypeFromSchema(*parameter.Schema)
				op.addDeps(subDeps...)
			} else {
				schema := spec.Schema{}
				schema.Typed(parameter.Type, parameter.Format)
				schema.Extensions = parameter.Extensions

				goType, subDeps = GetTypeFromSchema(schema)
				op.addDeps(subDeps...)
			}

			var tags []string

			if parameter.Type == "string" && goType != "string" {
				jsonTag = codegen.JoinWithComma(jsonTag, "string")
			}

			tags = append(tags, codegen.DeclTag("json", jsonTag), codegen.DeclTag("in", inTag))

			if fmt.Sprint(parameter.Default) != "<nil>" {
				tags = append(tags, codegen.DeclTag("default", fmt.Sprint(parameter.Default)))
			}

			if parameter.Extensions["x-go-validate"] != nil {
				tags = append(tags, codegen.DeclTag("validate", fmt.Sprint(parameter.Extensions["x-go-validate"])))
			}

			fields = append(fields, codegen.DeclField(
				fieldName,
				goType,
				tags,
				parameter.Description,
			))
		}

		return codegen.DeclType(op.ID+"Request", codegen.DeclStruct(fields))
	}

	return ""
}

func (op *OperationInfo) RenderRespDecl() string {
	schema := spec.Schema{}

	schema.Typed("object", "")

	schema.SetProperty("body", op.RespBody)

	goType, subDeps := ToGoType(op.ID+"Response", schema)

	op.addDeps(subDeps...)

	return goType
}

func (op *OperationInfo) RenderOperationMethod() string {
	if len(op.Parameters) > 0 {
		return codegen.TemplateRender(`{{ .ID }}(req {{ .ID }}Request) (resp {{ .ID }}Response, err error) {
			err = c.DoRequest("{{ .ServiceName }}.{{ .ID }}", "{{ .Method  }}", "{{ .Path }}", req, &resp)
			return
	}`)(op)
	}

	return codegen.TemplateRender(`{{ .ID }}() (resp {{ .ID }}Response, err error) {
		err = c.DoRequest("{{ .ServiceName }}.{{ .ID }}", "{{ .Method  }}", "{{ .Path }}", nil, &resp)
		return
	}`)(op)
}

func ToClient(baseClient string, pkgName string, swagger spec.Swagger) string {
	clientInfo := ClientInfo{
		BaseClient: baseClient,
		PkgName:    pkgName,
		Name:       codegen.ToUpperCamelCase(pkgName),
	}

	for path, pathItem := range swagger.Paths.Paths {
		if pathItem.PathItemProps.Get != nil {
			clientInfo.AddOperation("GET", path, pathItem.PathItemProps.Get, swagger.Responses)
		}
		if pathItem.PathItemProps.Post != nil {
			clientInfo.AddOperation("POST", path, pathItem.PathItemProps.Post, swagger.Responses)
		}
		if pathItem.PathItemProps.Put != nil {
			clientInfo.AddOperation("PUT", path, pathItem.PathItemProps.Put, swagger.Responses)
		}
		if pathItem.PathItemProps.Delete != nil {
			clientInfo.AddOperation("DELETE", path, pathItem.PathItemProps.Delete, swagger.Responses)
		}
		if pathItem.PathItemProps.Head != nil {
			clientInfo.AddOperation("HEAD", path, pathItem.PathItemProps.Head, swagger.Responses)
		}
		if pathItem.PathItemProps.Patch != nil {
			clientInfo.AddOperation("PATCH", path, pathItem.PathItemProps.Patch, swagger.Responses)
		}
		if pathItem.PathItemProps.Options != nil {
			clientInfo.AddOperation("OPTIONS", path, pathItem.PathItemProps.Options, swagger.Responses)
		}
	}

	return clientInfo.Render()
}
