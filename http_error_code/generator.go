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

func NewErrorGenerator(packagePath string, errorRegisterMethod string) *ErrorGenerator {
	prog := program.NewProgram(packagePath)

	return &ErrorGenerator{
		ErrorRegisterMethod: errorRegisterMethod,
		PackagePath:         packagePath,
		Program:             prog,
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
	ErrorRegisterMethod string
	PackagePath         string
	Program             *program.Program
	HttpErrors          map[*types.Package]map[string]HttpErrorValue
}

func CollectErrors(p *program.Program) map[*types.Package]map[string]HttpErrorValue {
	httpErrors := map[*types.Package]map[string]HttpErrorValue{}

	for _, pkgInfo := range p.AllPackages {
		for ident, obj := range pkgInfo.Defs {
			if constObj, ok := obj.(*types.Const); ok {
				if program.IsTypeName(obj.Type(), HttpErrorVarName) {
					name := constObj.Name()
					if name == "_" {
						continue
					}

					doc := program.GetTextFromCommentGroup(p.CommentGroupFor(ident))
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

			importedErrorType, method := program.ParsePkgExpose(g.ErrorRegisterMethod)

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
				ParseErrorRegister(sortedHttpErrorValues, method),
			}

			codegen.GenerateGoFile(
				codegen.JoinWithSlash(path, "errors.go"),
				strings.Join(blocks, "\n\n"),
			)
		}
	}
}

func ParseErrorRegister(httpErrorValues []HttpErrorValue, errorRegisterMethod string) string {
	codes := `func init () {
	`

	for _, httpErrorValue := range httpErrorValues {
		codes += errorRegisterMethod + `( ` + strings.Join([]string{
			strconv.Quote(httpErrorValue.Name),
			httpErrorValue.Code,
			strconv.Quote(httpErrorValue.Msg),
			strconv.Quote(httpErrorValue.Desc),
			strconv.FormatBool(httpErrorValue.CanBeErrTalk),
		}, ", ") + `)
		`
	}

	codes += `}`

	return codes
}
