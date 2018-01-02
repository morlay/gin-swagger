package swagger

import (
	"go/ast"
	"go/types"

	"github.com/morlay/gin-swagger/http_error_code"
	"github.com/morlay/gin-swagger/program"
)

func NewHttpErrorScanner() *HttpErrorScanner {
	return &HttpErrorScanner{}
}

type HttpErrorScanner struct {
	ErrorType            types.Type
	HttpErrors           map[*types.Package]map[string]http_error_code.HttpErrorValue
	funcUsesHttpErrors   map[*types.Func]map[string]http_error_code.HttpErrorValue
	funcMarkedHttpErrors map[*types.Func][]string
}

func (scanner *HttpErrorScanner) GetMarkedErrorsForFunc(funcType *types.Func) ([]string, bool) {
	httpErrors, ok := scanner.funcMarkedHttpErrors[funcType]
	return httpErrors, ok
}

func (scanner *HttpErrorScanner) GetErrorsInFunc(funcType *types.Func) (map[string]http_error_code.HttpErrorValue, bool) {
	httpErrors, ok := scanner.funcUsesHttpErrors[funcType]
	return httpErrors, ok
}

func (scanner *HttpErrorScanner) ForEachError(handler func(pkgDefHttpError *types.Package)) {
	for pkgDefHttpError := range scanner.HttpErrors {
		handler(pkgDefHttpError)
	}
}

func (scanner *HttpErrorScanner) Scan(prog *program.Program) {
	scanner.HttpErrors = http_error_code.CollectErrors(prog)

	scanner.funcUsesHttpErrors = map[*types.Func]map[string]http_error_code.HttpErrorValue{}
	scanner.funcMarkedHttpErrors = map[*types.Func][]string{}

	httpErrorMethods := map[*types.Func]*types.Signature{}

	for pkg, pkgInfo := range prog.AllPackages {
		for id, obj := range pkgInfo.Defs {
			if tpeFunc, ok := obj.(*types.Func); ok {
				if id.Obj != nil {
					if funcDecl, ok := id.Obj.Decl.(*ast.FuncDecl); ok {
						_, httpErrors := ParseHttpError(funcDecl.Doc.Text())
						if len(httpErrors) > 0 {
							scanner.funcMarkedHttpErrors[tpeFunc] = httpErrors
						}
					}
				}

				for pkgDefHttpError, httpErrorMap := range scanner.HttpErrors {
					if pkg == pkgDefHttpError || program.PkgContains(pkg.Imports(), pkgDefHttpError) {
						for id, obj := range pkgInfo.Uses {
							if tpeFunc.Scope() != nil && tpeFunc.Scope().Contains(id.Pos()) {
								if constObj, ok := obj.(*types.Const); ok {
									if program.IsTypeName(obj.Type(), http_error_code.HttpErrorVarName) {
										code := constObj.Val().String()
										if httpErrorValue, ok := httpErrorMap[code]; ok {
											if scanner.ErrorType == nil {
												httpErrorMethods = prog.MethodsOf(obj.Type())
												for funcType, method := range httpErrorMethods {
													if funcType.Name() == "ToError" {
														scanner.ErrorType = method.Results().At(0).Type()
													}
													if funcType.Name() == "StatusError" {
														scanner.ErrorType = method.Results().At(0).Type()
													}
												}
											}

											if httpErrorMethods[tpeFunc] == nil {
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
	}
}
