package swagger

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/logrusorgru/aurora"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/morlay/gin-swagger/codegen"
	"github.com/morlay/gin-swagger/http_error_code"
	"github.com/morlay/gin-swagger/program"
	"path"
)

func NewScanner(packagePath string) *Scanner {
	prog := program.NewProgram(packagePath)
	swag := NewSwagger()
	return &Scanner{
		Swagger: swag,
		Program: prog,
	}
}

type Scanner struct {
	GinPath            string
	Swagger            *Swagger
	Program            *program.Program
	httpErrors         map[*types.Package]map[string]http_error_code.HttpErrorValue
	funcUsesHttpErrors map[*types.Func]map[string]http_error_code.HttpErrorValue
}

func (scanner *Scanner) getRouterPrefixByIdent(id *ast.Ident) (string, []ast.Expr) {
	def := scanner.Program.ObjectOf(id)

	var prefix = ""

	args := []ast.Expr{}

	if def != nil {
		if assignStmt, ok := program.GetIdentDecl(id).(*ast.AssignStmt); ok {
			callExpr := assignStmt.Rhs[0].(*ast.CallExpr)
			if pointer, ok := def.Type().(*types.Pointer); ok {
				if !typeOfGinEngine(pointer) {
					if nextIdent, ok := callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident); ok {
						parentPrefix, parentArgs := scanner.getRouterPrefixByIdent(nextIdent)
						args = append(args, parentArgs...)
						if len(callExpr.Args) > 1 {
							args = append(args, callExpr.Args[1:]...)
						}
						return parentPrefix + getRouterPathByCallExpr(callExpr), args
					}
				}
			}
		}
	}

	return prefix, args
}

func (scanner *Scanner) getNodeDoc(node ast.Node) string {
	return program.GetTextFromCommentGroup(scanner.Program.CommentGroupFor(node))
}

func (scanner *Scanner) getEnums(doc string, node ast.Node) (enums []interface{}, enumLabels []string, enumVals []interface{}, otherDoc string) {
	var hasEnum bool
	otherDoc, hasEnum = ParseEnum(doc)
	if hasEnum {
		options := scanner.Program.GetEnumOptionsByType(node)
		for _, option := range options {
			enums = append(enums, option.Value)
			enumLabels = append(enumLabels, option.Label)
			enumVals = append(enumVals, option.V)
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

		if doc, fmtName = ParseStrfmt(doc); fmtName != "" {
			newSchema.Typed("string", fmtName)
			newSchema.WithDescription(doc)
			return newSchema
		}

		var enums []interface{}
		var enumLabels []string
		var enumVals []interface{}

		if enums, enumLabels, enumVals, doc = scanner.getEnums(doc, astType); len(enums) > 0 {
			if len(enums) == 2 && strings.ToUpper(enums[0].(string)) == "TRUE" && strings.ToUpper(enums[1].(string)) == "FALSE" {
				newSchema.Typed("boolean", "")
			} else {
				newSchema.WithEnum(enums...)
				if typeName, _, ok := GetSchemaTypeFromBasicType(reflect.TypeOf(enums[0]).Name()); ok {
					newSchema.Typed(typeName, "")
				}
				newSchema.AddExtension("x-enum-values", enums)
				newSchema.AddExtension("x-enum-labels", enumLabels)
				newSchema.AddExtension("x-enum-vals", enumVals)
				newSchema.WithDescription(doc)
			}
		}

		newSchema.WithDescription(doc)
	case *types.Basic:
		if typeName, format, ok := GetSchemaTypeFromBasicType(t.(*types.Basic).Name()); ok {
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
			log.Printf(aurora.Sprintf("\t Picking defination from %s\n", aurora.Brown(namedType)))
			s, ok := scanner.Swagger.AddDefinition(name, scanner.defineSchemaBy(namedType.Underlying()))
			if !ok {
				log.Printf(aurora.Sprintf(aurora.Red("\t\t `%s` already picked from `%s`"), name, namedType))
			}
			schema = *s
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
			structFieldTags := reflect.StructTag(structType.Tag(i))

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
						commonValidations := GetCommonValidations(validate)
						BindSchemaWithCommonValidations(&propSchema, commonValidations)
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

func (scanner *Scanner) getNonBodyParameter(name string, location string, tags reflect.StructTag, t types.Type) spec.Parameter {
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
			commonValidations := GetCommonValidations(validate)
			BindSchemaWithCommonValidations(&schema, commonValidations)
			schema.AddExtension("x-go-validate", validate)
		}

		BindItemsWithSchema(&items, schema)

		// todo support other collection format
		param.CollectionOf(&items, "csv")
	case *types.Basic, *types.Named:
		schema := scanner.getBasicSchemaFromType(t)

		if hasValidate {
			commonValidations := GetCommonValidations(validate)
			BindSchemaWithCommonValidations(&schema, commonValidations)
			schema.AddExtension("x-go-validate", validate)
		}

		BindParameterWithSchema(&param, schema)

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

func (scanner *Scanner) writeParameter(operation *spec.Operation, t types.Type) {
	if st, ok := t.(*types.Struct); ok {
		var structType = scanner.Program.WhereDecl(st).(*ast.StructType)

		for i := 0; i < st.NumFields(); i++ {
			var field = st.Field(i)
			var astField = structType.Fields.List[i]
			var structFieldTags = reflect.StructTag(st.Tag(i))
			var fieldType = field.Type()
			var fieldName = field.Name()

			if field.Anonymous() {
				scanner.writeParameter(operation, program.Indirect(fieldType))
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
		panic(fmt.Errorf("%s", "Param must be an struct"))
	}
}

func (scanner *Scanner) getStatusCodeFromExpr(expr ast.Expr) (int64, error) {
	constantValue := scanner.Program.ValueOf(expr)

	if constantValue == nil {
		return 0, fmt.Errorf("%s is not a constant value", expr)
	}

	return strconv.ParseInt(constantValue.String(), 10, 64)
}

func newOrMergeRespose(operation *spec.Operation, statusCode int) *spec.Response {
	var resp *spec.Response

	if operation.Responses != nil && operation.Responses.StatusCodeResponses != nil {
		r := operation.Responses.StatusCodeResponses[statusCode]
		resp = &r
	} else {
		resp = spec.NewResponse()
	}

	return resp
}

func (scanner *Scanner) writeResponse(operation *spec.Operation, ginContextCallExpr *ast.CallExpr, desc string) {
	args := ginContextCallExpr.Args

	statusCodeString, err := scanner.getStatusCodeFromExpr(args[0])

	if err == nil {
		statusCode := int(statusCodeString)
		resp := newOrMergeRespose(operation, statusCode)
		resp.WithDescription(resp.Description + desc)

		switch program.GetCallExprFunName(ginContextCallExpr) {
		// c.JSON(code int, obj interface{});
		case "JSON":
			if len(args) == 2 {
				tpe := scanner.Program.TypeOf(args[1])
				if !strings.Contains(tpe.String(), "untyped nil") {
					schema := scanner.defineSchemaBy(tpe)
					resp.WithSchema(&schema)
				}
				operation.Produces = []string{gin.MIMEJSON}
			}
		// c.HTML(code int, );
		// c.HTMLString(http.StatusOK, format, values)
		case "HTML", "HTMLString":
			operation.Produces = []string{gin.MIMEHTML}
		// c.String(http.StatusOK, format, values)
		case "String":
			schema := spec.Schema{}
			schema.Typed("string", "")
			resp.WithSchema(&schema)
		// c.Render(code init, )
		// c.Data(code init, )
		// c.Redirect(code init, )
		case "Render", "Data", "Redirect":
		}

		operation.RespondsWith(int(statusCode), resp)
	}
}

func (scanner *Scanner) writeResponseByHttpErrorValue(operation *spec.Operation, httpErrorValue http_error_code.HttpErrorValue, tpe types.Type) {
	statusCode := http_error_code.CodeToStatus(httpErrorValue.Code)

	resp := newOrMergeRespose(operation, statusCode)

	desc := resp.Description

	errDesc := `HttpError(` + httpErrorValue.Name + `,` + httpErrorValue.Code + `,` + strconv.Quote(httpErrorValue.Msg) + `,` + strconv.Quote(httpErrorValue.Desc) + `,` + fmt.Sprint(httpErrorValue.CanBeErrTalk) + `);`

	if strings.Index(desc, errDesc) == -1 {
		desc = strings.Join([]string{desc, errDesc}, "\n")
	}

	if !strings.Contains(tpe.String(), "untyped nil") {
		schema := scanner.defineSchemaBy(tpe)
		resp.WithSchema(&schema)
	}

	operation.Produces = []string{gin.MIMEJSON}
	resp.WithDescription(desc)
	operation.RespondsWith(statusCode, resp)
}

func (scanner *Scanner) pickOperationInfo(operation *spec.Operation, scope *types.Scope, scanned map[*types.Scope]bool) {
	scanned[scope] = true

	funType := scanner.Program.WitchFunc(scope.Pos())

	log.Printf("Picking operation from %s\n", aurora.Blue(funType.FullName()))

	for _, name := range scope.Names() {
		tpe := scope.Lookup(name).Type()
		// get parameters from type of var `req` or `request`;
		if name == "req" || name == "request" {
			if structTpe, ok := program.Indirect(tpe).(*types.Struct); ok {
				astStruct := scanner.Program.WhereDecl(tpe)
				log.Printf("\t Picking parameters from %s\n", aurora.Sprintf(aurora.Green("%s"), astStruct))
				scanner.writeParameter(operation, structTpe)
			} else {
				panic(fmt.Errorf(aurora.Sprintf(aurora.Red("request in %s should be a struct\n"))))
			}
		}
	}

	for id, obj := range scanner.Program.UsesInScope(scope) {
		switch obj.(type) {
		case *types.Func:
			tpeFunc := obj.(*types.Func)
			if isFuncOfGin(tpeFunc) {
				if callExpr := scanner.Program.CallFuncById(id); callExpr != nil {
					scanner.writeResponse(operation, callExpr, scanner.getNodeDoc(callExpr))
				}
			} else if !scanned[tpeFunc.Scope()] {
				if isFuncWithGinContext(tpeFunc) {
					scanner.pickOperationInfo(operation, tpeFunc.Scope(), scanned)
				} else if _, ok := scanner.funcUsesHttpErrors[tpeFunc]; ok {
					scanner.pickOperationInfo(operation, tpeFunc.Scope(), scanned)
				}
			}
		case *types.Const:
			if len(scanner.httpErrors) > 0 {
				constObj := obj.(*types.Const)
				if program.IsTypeName(obj.Type(), http_error_code.HttpErrorVarName) {
					if httpErrorValue, ok := scanner.httpErrors[obj.Pkg()][constObj.Val().String()]; ok {
						methods := scanner.Program.MethodsOf(obj.Type())
						for funcType, method := range methods {
							if funcType.Name() == "ToError" {
								scanner.writeResponseByHttpErrorValue(operation, httpErrorValue, method.Results().At(0).Type())
							}
						}
					}
				}
			}
		}
	}
}

func (scanner *Scanner) writeSummaryDesc(operation *spec.Operation, doc string) {
	summary, desc := parseCommentToSummaryDesc(doc)
	operation.WithSummary(summary)
	operation.WithDescription(desc)
}

func (scanner *Scanner) writeOperation(operation *spec.Operation, handlerFuncDecl *ast.FuncDecl) {
	scanned := map[*types.Scope]bool{}

	scope := scanner.Program.ScopeOf(handlerFuncDecl)
	scanner.pickOperationInfo(operation, scope, scanned)

	scanner.writeSummaryDesc(operation, handlerFuncDecl.Doc.Text())
}

func (scanner *Scanner) patchPathWithZero(swaggerPath string, operation *spec.Operation) string {
	r := regexp.MustCompile("/\\{([^/\\}]+)\\}")

	return r.ReplaceAllStringFunc(swaggerPath, func(str string) string {
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

		log.Printf(aurora.Sprintf(aurora.Red("`%s` without defining param `%s`, and use 0 instead;\n"), swaggerPath, name))

		return "/0"
	})
}

func patchOperationConsumes(operation *spec.Operation) {
	var isParameterHasBodySchema = false

	for _, parameter := range operation.Parameters {
		if parameter.In == "body" && parameter.Schema != nil {
			isParameterHasBodySchema = true
		}
	}

	if isParameterHasBodySchema {
		operation.WithConsumes(gin.MIMEJSON)
	}
}

func (scanner *Scanner) collectOperation(method string, ginPath string, handlerExprs []ast.Expr) {
	operation := new(spec.Operation)
	swaggerPath := convertGinPathToSwaggerPath(ginPath)

	log.Printf("%s %s\n", aurora.Red(method), aurora.Blue(ginPath))

	lastIdx := len(handlerExprs) - 1

	for idx, handlerExpr := range handlerExprs {
		var operationIdent *ast.Ident

		switch handlerExpr.(type) {
		case *ast.Ident:
			operationIdent = handlerExpr.(*ast.Ident)
		case *ast.SelectorExpr:
			operationIdent = handlerExpr.(*ast.SelectorExpr).Sel
		}

		ident := scanner.Program.IdentOf(scanner.Program.DefOf(operationIdent))

		if funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl); ok {
			scanner.writeOperation(operation, funcDecl)
			if idx == lastIdx {
				pkgInfo := scanner.Program.PackageInfoOf(funcDecl)
				operation.WithTags(pkgInfo.Pkg.Name())
				operation.WithID(operationIdent.String())
			}
		}
	}

	patchOperationConsumes(operation)

	scanner.Swagger.AddOperation(method, scanner.patchPathWithZero(swaggerPath, operation), operation)
}

func (scanner *Scanner) CollectHttpErrors() {
	scanner.httpErrors = http_error_code.CollectErrors(scanner.Program)
	scanner.funcUsesHttpErrors = map[*types.Func]map[string]http_error_code.HttpErrorValue{}

	for pkgDefHttpError, httpErrorMap := range scanner.httpErrors {
		for pkg, pkgInfo := range scanner.Program.AllPackages {
			if pkg == pkgDefHttpError || program.PkgContains(pkg.Imports(), pkgDefHttpError) {
				for _, obj := range pkgInfo.Defs {
					if tpeFunc, ok := obj.(*types.Func); ok {
						for _, obj := range pkgInfo.Uses {
							if constObj, ok := obj.(*types.Const); ok {
								if program.IsTypeName(obj.Type(), http_error_code.HttpErrorVarName) {
									code := constObj.Val().String()
									if httpErrorValue, ok := httpErrorMap[code]; ok {
										if scanner.funcUsesHttpErrors[tpeFunc] == nil {
											scanner.funcUsesHttpErrors[tpeFunc] = map[string]http_error_code.HttpErrorValue{}
										}
										scanner.funcUsesHttpErrors[tpeFunc][code] = httpErrorValue
									}
								}
							}
						}
					}
				}
			}
		}

	}
}

func (scanner *Scanner) Scan() {
	scanner.CollectHttpErrors()

	for pkg, pkgInfo := range scanner.Program.AllPackages {
		if hasImportedGin(pkg.Imports()) {
			for selectorExpr, selection := range pkgInfo.Info.Selections {
				if pointer, ok := selection.Recv().(*types.Pointer); ok {
					if typeOfGinEngine(pointer) || typeOfGinRouterGroup(pointer) {
						if isGinMethod(selectorExpr.Sel.Name) {
							if callExpr := program.FindCallExprByFunc(pkgInfo.Info, selectorExpr); callExpr != nil {
								method := selectorExpr.Sel.Name
								prefix := ""
								args := []ast.Expr{}

								if typeOfGinRouterGroup(pointer) {
									prefix, args = scanner.getRouterPrefixByIdent(selectorExpr.X.(*ast.Ident))
								}

								ginPath := path.Join(prefix, getRouterPathByCallExpr(callExpr))

								args = append(args, callExpr.Args[1:]...)

								scanner.collectOperation(method, ginPath, args)
							}
						}
					}
				}
			}
		}
	}
}

func (scanner *Scanner) Output(path string) {
	scanner.Scan()
	codegen.WriteJSONFile(path, scanner.Swagger)
}
