package scanner

import (
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/morlay/gin-swagger/program"
	"github.com/morlay/gin-swagger/swagger"
	"gopkg.in/gin-gonic/gin.v1"
)

type ScannerOpts struct {
	PackagePath string
}

func NewScanner(opts *ScannerOpts) *Scanner {
	prog := program.NewProgram(opts.PackagePath)
	swag := swagger.NewSwagger()
	return &Scanner{
		Swagger: swag,
		Program: prog,
	}
}

type Scanner struct {
	GinPath string
	Swagger *swagger.Swagger
	Program *program.Program
}

func (scanner *Scanner) packageOfGin(packagePath string) bool {
	return strings.Contains(packagePath, "gin")
}

func (scanner *Scanner) typeOfGinEngine(pointer *types.Pointer) bool {
	packagePath := pointer.String()
	return scanner.packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "Engine"
}

func (scanner *Scanner) typeOfGinRouterGroup(pointer *types.Pointer) bool {
	packagePath := pointer.String()
	return scanner.packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "RouterGroup"
}

func (scanner *Scanner) typeOfGinContext(pointer *types.Pointer) bool {
	packagePath := pointer.String()
	return scanner.packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "Context"
}

func (scanner *Scanner) getRouterPrefixByIdent(id *ast.Ident) string {
	def := scanner.Program.ObjectOf(id)

	var prefix = ""

	if def != nil {
		if assignStmt, ok := program.GetIdentDecl(id).(*ast.AssignStmt); ok {
			callExpr := assignStmt.Rhs[0].(*ast.CallExpr)
			if pointer, ok := def.Type().(*types.Pointer); ok {
				if !scanner.typeOfGinEngine(pointer) {
					if nextIdent, ok := callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident); ok {
						return scanner.getRouterPrefixByIdent(nextIdent) + getRouterPathByCallExpr(callExpr)
					}
				}
			}
		}
	}

	return prefix
}

func (scanner *Scanner) getNodeDoc(node ast.Node) string {
	return program.GetTextFromCommentGroup(scanner.Program.CommentGroupFor(node))
}

func (scanner *Scanner) getStrFmt(doc string) (fmtName string, otherDoc string) {
	otherDoc, fmtName = swagger.ParseStrfmt(doc)
	return
}

func (scanner *Scanner) getEnums(doc string, node ast.Node) (enums []interface{}, enumLabels []string, otherDoc string) {
	var hasEnum bool
	otherDoc, hasEnum = swagger.ParseEnum(doc)
	if hasEnum {
		options := scanner.Program.GetEnumOptionsByType(node)
		for _, option := range options {
			enums = append(enums, option.Value)
			enumLabels = append(enumLabels, option.Label)
		}
	}
	return
}

func (scanner *Scanner) getBasicSchemaFromType(t types.Type) spec.Schema {
	var newSchema = spec.Schema{}

	switch t.(type) {
	case *types.Named:
		namedType := t.(*types.Named)
		astType := scanner.Program.WhereDecl(namedType)
		newSchema = scanner.getBasicSchemaFromType(namedType.Underlying())

		newSchema.AddExtension("x-go-named", namedType.String())

		var fmtName string
		var doc string

		doc = scanner.getNodeDoc(astType)

		if fmtName, doc = scanner.getStrFmt(doc); fmtName != "" {
			newSchema.Typed("string", fmtName)
			newSchema.WithDescription(doc)
			return newSchema
		}

		var enums []interface{}
		var enumLabels []string

		if enums, enumLabels, doc = scanner.getEnums(doc, astType); len(enums) > 0 {
			newSchema.WithEnum(enums...)
			if typeName, _, ok := swagger.GetSchemaTypeFromBasicType(reflect.TypeOf(enums[0]).Name()); ok {
				newSchema.Typed(typeName, "")
			}
			newSchema.AddExtension("x-enum-values", enums)
			newSchema.AddExtension("x-enum-labels", enumLabels)
			newSchema.WithDescription(doc)
		}

		newSchema.WithDescription(doc)
	case *types.Basic:
		if typeName, format, ok := swagger.GetSchemaTypeFromBasicType(t.(*types.Basic).Name()); ok {
			newSchema.Typed(typeName, format)
		}
	}
	return newSchema
}

func (scanner *Scanner) defineSchemaBy(tpe types.Type) spec.Schema {
	schema := spec.Schema{}

	switch tpe.(type) {
	case *types.Basic:
		schema = scanner.getBasicSchemaFromType(tpe)
	case *types.Named:
		namedType := tpe.(*types.Named)
		schema = scanner.getBasicSchemaFromType(tpe)

		if len(schema.Type) == 0 {
			name := getExportedNameOfPackage(namedType.String())
			schema = scanner.Swagger.AddDefinition(name, scanner.defineSchemaBy(namedType.Underlying()))
		}
	case *types.Pointer:
		schema = scanner.defineSchemaBy(tpe.(*types.Pointer).Elem())
	case *types.Map:
		propSchema := scanner.defineSchemaBy(tpe.(*types.Map).Elem())
		schema = *spec.MapProperty(&propSchema)
	case *types.Slice:
		itemsSchema := scanner.defineSchemaBy(tpe.(*types.Slice).Elem())
		schema = *spec.ArrayProperty(&itemsSchema)
	case *types.Struct:
		var structType = tpe.(*types.Struct)
		var structTypeAst = scanner.Program.WhereDecl(structType).(*ast.StructType)

		var structSchema = spec.Schema{}
		var schemas []spec.Schema

		structSchema.Typed("object", "")
		structSchema.WithDescription(scanner.getNodeDoc(structTypeAst))

		for i := 0; i < structType.NumFields(); i++ {
			field := structType.Field(i)
			fieldAst := structTypeAst.Fields.List[i]
			structFieldType := field.Type()
			structFieldTags := program.StructTag(structType.Tag(i))

			if field.Anonymous() {
				schemas = append(schemas, scanner.defineSchemaBy(structFieldType))
			} else {
				name, flags := getJSONNameAndFlags(structFieldTags.Get("json"))

				if name == "-" {
					continue
				}

				defaultValue, hasDefault := structFieldTags.Lookup("default")
				validate, hasValidate := structFieldTags.Lookup("validate")

				if name == "" {
					panic(fmt.Errorf("missing tag json for %s.%s\b", structType, field.Name()))
				}

				propSchema := scanner.defineSchemaBy(structFieldType)

				if len(flags) == 1 {
					propSchema.Typed("string", propSchema.Format)
				}

				propSchema.WithDescription(scanner.getNodeDoc(fieldAst))

				if hasDefault {
					propSchema.WithDefault(defaultValue)
				} else {
					structSchema.AddRequired(name)
				}

				if hasValidate {
					propSchema.WithDefault(defaultValue)
					propSchema.AddExtension("x-go-validate", validate)

					if hasValidate {
						commonValidations := swagger.GetCommonValidations(validate)
						swagger.BindSchemaWithCommonValidations(&propSchema, commonValidations)
					}
				}

				propSchema.AddExtension("x-go-name", field.Name())
				structSchema.SetProperty(name, propSchema)
			}

		}

		if len(schemas) > 0 {
			schemas = append(schemas, structSchema)
			schema.WithAllOf(schemas...)
		} else {
			schema = structSchema
		}
	}

	return schema
}

func (scanner *Scanner) getBodyParameter(t types.Type) spec.Parameter {
	schema := scanner.defineSchemaBy(t)
	return *spec.BodyParam("body", &schema)
}

func (scanner *Scanner) getNonBodyParameter(name string, location string, tags program.StructTag, t types.Type) spec.Parameter {
	param := spec.Parameter{}

	defaultValue, hasDefault := tags.Lookup("default")
	validate, hasValidate := tags.Lookup("validate")

	switch t.(type) {
	case *types.Pointer:
		param = scanner.getNonBodyParameter(name, location, tags, t.(*types.Pointer).Elem())
		return param
	case *types.Slice:
		sliceElem := t.(*types.Slice).Elem()
		var schema spec.Schema
		items := spec.Items{}

		switch sliceElem.(type) {
		case *types.Pointer:
			schema = scanner.getBasicSchemaFromType(sliceElem.(*types.Pointer).Elem())
		case *types.Named, *types.Basic:
			schema = scanner.getBasicSchemaFromType(sliceElem)
		}

		if hasValidate {
			commonValidations := swagger.GetCommonValidations(validate)
			swagger.BindSchemaWithCommonValidations(&schema, commonValidations)
			schema.AddExtension("x-go-validate", validate)
		}

		swagger.BindItemsWithSchema(&items, schema)

		// todo support other collection format
		param.CollectionOf(&items, "csv")
	case *types.Basic, *types.Named:
		schema := scanner.getBasicSchemaFromType(t)

		if hasValidate {
			commonValidations := swagger.GetCommonValidations(validate)
			swagger.BindSchemaWithCommonValidations(&schema, commonValidations)
			schema.AddExtension("x-go-validate", validate)
		}

		swagger.BindParameterWithSchema(&param, schema)

	}

	if !hasDefault {
		param.AsRequired()
	} else {
		param.WithDefault(defaultValue)
	}

	param.WithLocation(location)
	param.Named(name)

	return param
}

func (scanner *Scanner) bindParamBy(t types.Type, operation *spec.Operation) {
	if st, ok := t.(*types.Struct); ok {
		var structType = scanner.Program.WhereDecl(st).(*ast.StructType)

		for i := 0; i < st.NumFields(); i++ {
			var field = st.Field(i)
			var astField = structType.Fields.List[i]
			var structFieldTags = program.StructTag(st.Tag(i))
			var fieldType = field.Type()
			var fieldName = field.Name()

			if field.Anonymous() {
				scanner.bindParamBy(indirect(fieldType), operation)
			} else {
				location := structFieldTags.Get("in")
				name, flags := getJSONNameAndFlags(structFieldTags.Get("json"))

				if location == "" {
					if fieldName == "Body" {
						location = "body"
					} else {
						panic(fmt.Errorf("missing tag `in` for %s.%s", st.String(), fieldName))
					}
				}

				if name == "" {
					if fieldName == "Body" {
						name = "body"
					} else {
						panic(fmt.Errorf("missing tag `json` for %s", fieldName))
					}
				}

				var param spec.Parameter

				if location == "body" {
					param = scanner.getBodyParameter(fieldType)
				} else {
					param = scanner.getNonBodyParameter(name, location, structFieldTags, fieldType)
					if len(flags) > 0 {
						// todo check other flags;
						param.Typed("string", param.Format)
					}
				}

				param.AddExtension("x-go-name", field.Name())
				param.WithDescription(scanner.getNodeDoc(astField))
				operation.AddParam(&param)
			}
		}
	} else {
		fmt.Errorf("%s", "Param must be an struct")
	}
}

func (scanner *Scanner) getStatusCodeFromExpr(expr ast.Expr) (float64, error) {
	return strconv.ParseInt(scanner.Program.ValueOf(expr).String(), 10, 64)
}

func (scanner *Scanner) addResponse(ginContextCallExpr *ast.CallExpr, desc string, operation *spec.Operation) {
	response := spec.NewResponse()
	args := ginContextCallExpr.Args

	response.WithDescription(desc)

	switch program.GetCallExprFunName(ginContextCallExpr) {
	// c.JSON(code int, obj interface{});
	case "JSON":
		statusCode, _ := scanner.getStatusCodeFromExpr(args[0]);
		tpe := scanner.Program.TypeOf(args[1])

		if !strings.Contains(tpe.String(), "untyped nil") {
			schema := scanner.defineSchemaBy(tpe)
			response.WithSchema(&schema)
		}

		operation.RespondsWith(int(statusCode), response)
		operation.WithProduces(gin.MIMEJSON)
	// c.HTML(code int, );
	case "HTML":
		statusCode, _ := scanner.getStatusCodeFromExpr(args[0]);
		operation.RespondsWith(int(statusCode), response)
		operation.WithProduces(gin.MIMEHTML)
	// c.String(http.StatusOK, format, values)
	case "String":
		statusCode, _ := scanner.getStatusCodeFromExpr(args[0]);
		schema := spec.Schema{}
		schema.Typed("string", "")
		response.WithSchema(&schema)
		operation.RespondsWith(int(statusCode), response)
	// c.Render(code init, )
	// c.Data(code init, )
	// c.Redirect(code init, )
	case "Render", "Data", "Redirect":
		statusCode, _ := scanner.getStatusCodeFromExpr(args[0]);
		operation.RespondsWith(int(statusCode), response)
	}
}

func (scanner *Scanner) getOperation(handlerFuncDecl *ast.FuncDecl) (operation *spec.Operation) {
	operation = new(spec.Operation)

	pkg := scanner.Program.PackageOf(handlerFuncDecl)
	scope := scanner.Program.ScopeOf(handlerFuncDecl)
	summary, desc := parseCommentToSummaryDesc(handlerFuncDecl.Doc.Text())

	operation.WithSummary(summary)
	operation.WithDescription(desc)
	operation.WithTags(pkg.Name())

	fmt.Printf("Get operation from `%s.%s`\n", pkg.Name(), handlerFuncDecl.Name.String())

	for _, name := range scope.Names() {
		tpe := scope.Lookup(name).Type()
		// get parameters from type of var `req` or `request`;
		if name == "req" || name == "request" {
			scanner.bindParamBy(tpe.Underlying(), operation)
		}

		// get response from method gin.Context
		if pointer, ok := tpe.(*types.Pointer); ok {
			if scanner.typeOfGinContext(pointer) {
				for _, pkgInfo := range scanner.Program.AllPackages {
					program.PickSelectionBy(pkgInfo.Info, func(selectorExpr *ast.SelectorExpr, selection *types.Selection) bool {
						if selection.Recv() == tpe {
							if callExpr := program.FindCallExprByFunc(pkgInfo.Info, selectorExpr); callExpr != nil {
								scanner.addResponse(callExpr, scanner.getNodeDoc(callExpr), operation)
							}
						}
						return false
					})
				}
			}
		}

	}

	var isParameterHasBodySchema = false

	for _, parameter := range operation.Parameters {
		if parameter.In == "body" && parameter.Schema != nil {
			isParameterHasBodySchema = true
		}
	}

	if isParameterHasBodySchema {
		operation.WithConsumes(gin.MIMEJSON)
	}

	return
}

func (scanner *Scanner) HasImportedGin(packages []*types.Package) bool {
	for _, pkg := range packages {
		if pkg.Name() == "gin" && scanner.packageOfGin(pkg.Path()) {
			return true
		}
	}
	return false
}

func (scanner *Scanner) collectOperationByCallExpr(callExpr *ast.CallExpr, prefix string) {
	method := program.GetCallExprFunName(callExpr)

	if isGinMethod(method) {
		args := callExpr.Args
		lastArg := args[len(args) - 1]

		var id string
		var swaggerPath string
		var operation *spec.Operation

		ginPath := path.Join(prefix, getRouterPathByCallExpr(callExpr))
		swaggerPath = convertGinPathToSwaggerPath(ginPath)

		var operationIdent *ast.Ident

		switch lastArg.(type) {
		case *ast.Ident:
			operationIdent = lastArg.(*ast.Ident)
		case *ast.SelectorExpr:
			operationIdent = lastArg.(*ast.SelectorExpr).Sel
		}

		id = operationIdent.String()

		ident := scanner.Program.IdentOf(scanner.Program.DefOf(operationIdent))

		if funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl); ok {
			operation = scanner.getOperation(funcDecl)
			operation.WithID(id)

			r := regexp.MustCompile("/\\{([^/\\}]+)\\}")

			fixedPath := r.ReplaceAllStringFunc(swaggerPath, func(str string) string {
				name := r.FindAllStringSubmatch(str, -1)[0][1]

				var isParameterDefined = false

				for _, parameter := range operation.Parameters {
					if parameter.In == "path" && parameter.Name == name {
						isParameterDefined = true
					}
				}

				if isParameterDefined {
					return str
				}

				fmt.Printf("%s without defining param `%s`, and use 0 instead;\n", swaggerPath, name)

				return "/0"
			})

			scanner.Swagger.AddOperation(method, fixedPath, operation)
		}
	}
}

func (scanner *Scanner) Scan() {
	for pkg, pkgInfo := range scanner.Program.AllPackages {
		if scanner.HasImportedGin(pkg.Imports()) {
			program.PickSelectionBy(pkgInfo.Info, func(selectorExpr *ast.SelectorExpr, selection *types.Selection) bool {
				if pointer, ok := selection.Recv().(*types.Pointer); ok {
					if scanner.typeOfGinEngine(pointer) || scanner.typeOfGinRouterGroup(pointer) {
						if isGinMethod(selectorExpr.Sel.Name) {
							if callExpr := program.FindCallExprByFunc(pkgInfo.Info, selectorExpr); callExpr != nil {
								var prefix = ""

								if scanner.typeOfGinRouterGroup(pointer) {
									prefix = scanner.getRouterPrefixByIdent(selectorExpr.X.(*ast.Ident))
								}

								scanner.collectOperationByCallExpr(callExpr, prefix)
							}
							return false
						}
					}
				}
				return false
			})
		}
	}
}
