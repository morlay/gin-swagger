package http_error_code

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sort"

	"go/build"

	"github.com/morlay/gin-swagger/codegen"
	"github.com/morlay/gin-swagger/program"
)

func NewErrorGenerator(packagePath string, errorType string) *ErrorGenerator {
	prog := program.NewProgram(packagePath)

	return &ErrorGenerator{
		ErrorType:   errorType,
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

func (h *HttpErrorValue) ToStatus() int {
	return CodeToStatus(h.Code)
}

func (h *HttpErrorValue) ToDesc() string {
	return `@httpError(` + h.Code + `,` + h.Name + `,` + strconv.Quote(h.Msg) + `,` + strconv.Quote(h.Desc) + `,` + fmt.Sprint(h.CanBeErrTalk) + `);`
}

type ByHttpErrorValue []HttpErrorValue

func (a ByHttpErrorValue) Len() int {
	return len(a)
}
func (a ByHttpErrorValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByHttpErrorValue) Less(i, j int) bool {
	return a[i].Name < a[j].Name
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
				if program.IsTypeName(obj.Type(), HttpErrorVarName) {
					doc := program.GetTextFromCommentGroup(p.CommentGroupFor(ident))
					name := constObj.Name()
					code := constObj.Val().String()
					msg, desc, canBeErrTalk := ParseHttpCodeDesc(doc)
					pkg := constObj.Pkg()

					if httpErrors[pkg] == nil {
						httpErrors[pkg] = map[string]HttpErrorValue{}
					}

					httpErrors[pkg][code] = HttpErrorValue{
						Name:         name,
						Code:         code,
						Msg:          msg,
						Desc:         desc,
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
			sortedHttpErrorValues := []HttpErrorValue{}

			for _, value := range httpErrorValues {
				sortedHttpErrorValues = append(sortedHttpErrorValues, value)
			}

			sort.Sort(ByHttpErrorValue(sortedHttpErrorValues))

			p, _ := build.Import(pkg.Path(), "", build.FindOnly)
			path, _ := filepath.Rel(cwd, p.Dir)

			importedErrorType, errorType := program.ParsePkgExpose(g.ErrorType)

			var imports = []string{
				"strconv",
				"fmt",
			}

			if importedErrorType != "" {
				imports = append(imports, importedErrorType)
			}

			blocks := []string{
				codegen.DeclPackage(pkg.Name()),
				codegen.DeclImports(imports...),
				ParseOthers(errorType),
				ParseKeyParser(sortedHttpErrorValues),
				ParseMsgParser(sortedHttpErrorValues),
				ParseDescParser(sortedHttpErrorValues),
				ParseErrorTalkParser(sortedHttpErrorValues),
			}

			codegen.GenerateGoFile(
				codegen.JoinWithSlash(path, "errors.go"),
				strings.Join(blocks, "\n\n"),
			)
		}
	}
}

func ParseMsgParser(httpErrorValues []HttpErrorValue) string {
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

func ParseKeyParser(httpErrorValues []HttpErrorValue) string {
	firstLine := `func (c ` + HttpErrorVarName + `) Key() string {
	switch (c) {`

	lines := []string{firstLine}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.Quote(httpErrorValue.Name)))
	}

	lines = append(lines, "}")
	lines = append(lines, codegen.DeclReturn(strconv.Quote("")))
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

func ParseDescParser(httpErrorValues []HttpErrorValue) string {
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

func ParseErrorTalkParser(httpErrorValues []HttpErrorValue) string {
	lines := []string{`func (c ` + HttpErrorVarName + `) CanBeErrTalk() bool {
	switch (c) {`}

	for _, httpErrorValue := range httpErrorValues {
		lines = append(lines, codegen.DeclCase(httpErrorValue.Name))
		lines = append(lines, codegen.DeclReturn(strconv.FormatBool(httpErrorValue.CanBeErrTalk)))
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
		Key:            c.Key(),
		Code:           c.Code(),
		Msg:            c.Msg(),
		Desc:           c.Desc(),
		CanBeErrorTalk: c.CanBeErrTalk(),
	}
}
`
}
