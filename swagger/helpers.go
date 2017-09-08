package swagger

import (
	"regexp"
	"strings"

	"go/ast"

	"github.com/morlay/gin-swagger/program"
)

func isGinMethod(method string) bool {
	var ginMethods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"PATCH":   true,
		"HEAD":    true,
		"DELETE":  true,
		"OPTIONS": true,
	}

	return ginMethods[method]
}

func isGinGroup(method string) bool {
	return method == "Group"
}

func getJSONNameAndFlags(tagValue string) (name string, flags map[string]bool) {
	values := strings.SplitN(tagValue, ",", -1)
	for _, flag := range values[1:] {
		if flags == nil {
			flags = map[string]bool{}
		}
		flags[flag] = true
	}
	name = values[0]
	return
}

func parseCommentToSummaryDesc(str string) (string, string) {
	lines := strings.SplitN(str, "\n", -1)
	return lines[0], strings.TrimSpace(strings.Join(lines[1:], "\n"))
}

func getExportedNameOfPackage(path string) string {
	var parts = strings.Split(path, ".")
	return parts[len(parts)-1]
}

func getRouterPathByCallExpr(callExpr *ast.CallExpr) string {
	return program.GetBasicLitValue(callExpr.Args[0].(*ast.BasicLit)).(string)
}

func convertGinPathToSwaggerPath(str string) string {
	r := regexp.MustCompile("/:([^/]+)")
	result := r.ReplaceAllString(str, "/{$1}")
	return result
}
