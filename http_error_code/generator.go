package http_error_code

import (
	"go/types"
	"path/filepath"
	"strings"
	"os"
	"strconv"
	"fmt"

	"github.com/morlay/gin-swagger/program"
	"github.com/morlay/gin-swagger/codegen"
)

func NewErrorGenerator(packagePath string, errorType string) *ErrorGenerator {
	prog := program.NewProgram(packagePath)

	return &ErrorGenerator{
		ErrorType: errorType,
		PackagePath: packagePath,
		Program:     prog,
	}
}

const HttpErrorVarName = "HttpErrorCode"

type HttpErrorValue struct {
	Name         string
	Code         string
	Msg          string
	Desc         string
	CanBeErrTalk bool
}

type ErrorGenerator struct {
	ErrorType   string
	PackagePath string
	Program     *program.Program
	HttpErrors  map[*types.Package][]HttpErrorValue
}

func (g *ErrorGenerator) addHttpError(pkg *types.Package, name string, code string, msg string, desc string, canBeErrTalk bool) {
	if g.HttpErrors == nil {
		g.HttpErrors = map[*types.Package][]HttpErrorValue{}
	}

	if (g.HttpErrors[pkg] == nil) {
		g.HttpErrors[pkg] = []HttpErrorValue{}
	}

	g.HttpErrors[pkg] = append(g.HttpErrors[pkg], HttpErrorValue{
		Name: name,
		Code: code,
		Msg: msg,
		Desc: desc,
		CanBeErrTalk: canBeErrTalk,
	})
}

func (g *ErrorGenerator) Scan() {
	for _, pkgInfo := range g.Program.AllPackages {
		if program.IsSubPackageOf(g.PackagePath, pkgInfo.Pkg.Path()) {
			for ident, obj := range pkgInfo.Defs {
				if constObj, ok := obj.(*types.Const); ok {
					if (program.IsTypeName(obj.Type(), HttpErrorVarName)) {
						doc := program.GetTextFromCommentGroup(g.Program.CommentGroupFor(ident))
						name := constObj.Name();
						code := constObj.Val().String();
						msg, desc, canBeErrTalk := ParseHttpCodeDesc(doc)

						g.addHttpError(
							constObj.Pkg(),
							name,
							code,
							msg,
							desc,
							canBeErrTalk,
						)
					}
				}
			}
		}
	}

}

func (g *ErrorGenerator) Output() {
	g.Scan()

	cwd, _ := os.Getwd()

	for pkg, httpErrorValues := range g.HttpErrors {
		path, _ := filepath.Rel(cwd, filepath.Join(os.Getenv("GOPATH"), "src", pkg.Path()))

		importedErrorType, errorType := program.ParsePkgExpose(g.ErrorType)

		blocks := []string{
			codegen.DeclPackage(pkg.Name()),
			codegen.DeclImports("strconv", "fmt", importedErrorType),
			ParseOthers(errorType),
			ParseMsgParser(httpErrorValues),
			ParseDescParser(httpErrorValues),
			ParseErrorTalkParser(httpErrorValues),
		}

		codegen.WriteGoFile(
			codegen.JoinWithSlash(
				path,
				"generated_errors.go",
			),
			strings.Join(blocks, "\n\n"),
		)
	}
}

func ParseMsgParser(httpErrorValues []HttpErrorValue) string {
	lines := []string{`func (httpErrorCode HttpErrorCode) Msg() string {
	switch (httpErrorCode) {`}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.Quote(httpErrorValue.Msg)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(strconv.Quote("")))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseDescParser(httpErrorValues []HttpErrorValue) string {
	lines := []string{`func (httpErrorCode HttpErrorCode) Desc() string {
	switch (httpErrorCode) {`}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.Quote(httpErrorValue.Desc)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(strconv.Quote("")))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseErrorTalkParser(httpErrorValues []HttpErrorValue) string {
	lines := []string{`func (httpErrorCode HttpErrorCode) CanBeErrTalk() bool {
	switch (httpErrorCode) {`}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(fmt.Sprintln(httpErrorValue.CanBeErrTalk)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn("false"))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseOthers(errorTypeSelector string) string {
	return `
func (httpErrorCode HttpErrorCode) Code() int32 {
	return int32(httpErrorCode)
}

func (httpErrorCode HttpErrorCode) Status() int {
	status, _ := strconv.Atoi(fmt.Sprintln(httpErrorCode)[:3])
	return status
}

func (httpErrorCode HttpErrorCode) ToError() *` + errorTypeSelector + `{
	return &` + errorTypeSelector + `{
		Code:           httpErrorCode.Code(),
		Msg:            httpErrorCode.Msg(),
		Desc:           httpErrorCode.Desc(),
		CanBeErrorTalk: httpErrorCode.CanBeErrTalk(),
	}
}

func (httpErrorCode HttpErrorCode) ToResp() (int, *` + errorTypeSelector + `) {
	return httpErrorCode.Status(), httpErrorCode.ToError()
}
`
}