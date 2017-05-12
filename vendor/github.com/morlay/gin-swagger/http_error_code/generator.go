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

var HttpErrorVarName = "HttpErrorCode"

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
	HttpErrors  map[*types.Package]map[string]HttpErrorValue
}

func CollectErrors(p *program.Program) map[*types.Package]map[string]HttpErrorValue {
	httpErrors := map[*types.Package]map[string]HttpErrorValue{}

	for _, pkgInfo := range p.AllPackages {
		for ident, obj := range pkgInfo.Defs {
			if constObj, ok := obj.(*types.Const); ok {
				if (program.IsTypeName(obj.Type(), HttpErrorVarName)) {
					doc := program.GetTextFromCommentGroup(p.CommentGroupFor(ident))
					name := constObj.Name();
					code := constObj.Val().String();
					msg, desc, canBeErrTalk := ParseHttpCodeDesc(doc)
					pkg := constObj.Pkg();

					if (httpErrors[pkg] == nil) {
						httpErrors[pkg] = map[string]HttpErrorValue{}
					}

					httpErrors[pkg][code] = HttpErrorValue{
						Name: name,
						Code: code,
						Msg: msg,
						Desc: desc,
						CanBeErrTalk: canBeErrTalk,
					}
				}
			}
		}

	}

	return httpErrors
}

func (g *ErrorGenerator) Output() {
	g.HttpErrors = CollectErrors(g.Program)

	cwd, _ := os.Getwd()

	for pkg, httpErrorValues := range g.HttpErrors {
		if program.IsSubPackageOf(g.PackagePath, pkg.Path()) {
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
}

func ParseMsgParser(httpErrorValues map[string]HttpErrorValue) string {
	firstLine := `func (c ` + HttpErrorVarName + `) Msg() string {
	switch (c) {`

	lines := []string{firstLine}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.Quote(httpErrorValue.Msg)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(strconv.Quote("")))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseDescParser(httpErrorValues map[string]HttpErrorValue) string {
	lines := []string{`func (c ` + HttpErrorVarName + `) Desc() string {
	switch (c) {`}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.Quote(httpErrorValue.Desc)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(strconv.Quote("")))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseErrorTalkParser(httpErrorValues map[string]HttpErrorValue) string {
	lines := []string{`func (c ` + HttpErrorVarName + `) CanBeErrTalk() bool {
	switch (c) {`}

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
func (c ` + HttpErrorVarName + `) Code() int32 {
	return int32(c)
}

func (c ` + HttpErrorVarName + `) Status() int {
	status, _ := strconv.Atoi(fmt.Sprintln(c)[:3])
	return status
}

func (c ` + HttpErrorVarName + `) ToError() *` + errorTypeSelector + `{
	return &` + errorTypeSelector + `{
		Code:           c.Code(),
		Msg:            c.Msg(),
		Desc:           c.Desc(),
		CanBeErrorTalk: c.CanBeErrTalk(),
	}
}

func (c ` + HttpErrorVarName + `) ToResp() (int, *` + errorTypeSelector + `) {
	return c.Status(), c.ToError()
}
`
}