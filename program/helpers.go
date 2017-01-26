package program

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"strings"
	"unicode"
)

func GetBasicLitValue(basicLit *ast.BasicLit) interface{} {
	switch basicLit.Kind.String() {
	case "INT":
		if result, err := strconv.ParseInt(basicLit.Value, 10, 64); err == nil {
			return result
		}
	case "FLOAT":
		if result, err := strconv.ParseFloat(basicLit.Value, 64); err == nil {
			return result
		}
	default:
		return strings.Trim(basicLit.Value, "\"")
	}
	return nil
}

func GetIdentDecl(ident *ast.Ident) interface{} {
	if ident.Obj == nil {
		fmt.Errorf("Ident %s has empty obj", ident.Name)
	}
	return ident.Obj.Decl
}

func GetCallExprFunName(callExpr *ast.CallExpr) string {
	return callExpr.Fun.(*ast.SelectorExpr).Sel.Name
}

func FindExprBy(info types.Info, pick func(expr ast.Expr) bool) ast.Expr {
	for expr := range info.Types {
		if pick(expr) {
			return expr
		}
	}
	return nil
}

func FindCallExprByFunc(info types.Info, funcExpr ast.Expr) *ast.CallExpr {
	callExpr := FindExprBy(info, func(expr ast.Expr) bool {
		if callExpr, ok := expr.(*ast.CallExpr); ok {
			return callExpr.Fun == funcExpr
		}
		return false
	})
	if callExpr != nil {
		return callExpr.(*ast.CallExpr)
	}
	return nil
}

func PickSelectionBy(info types.Info, picker func(selectorExpr *ast.SelectorExpr, selection *types.Selection) bool) map[*ast.SelectorExpr]*types.Selection {
	var selections = make(map[*ast.SelectorExpr]*types.Selection)

	for selectorExpr, selection := range info.Selections {
		if picker(selectorExpr, selection) {
			selections[selectorExpr] = selection
		}
	}

	return selections
}

func GetTextFromCommentGroup(commentGroup []*ast.CommentGroup) string {
	var text = ""

	for _, comment := range commentGroup {
		text = text + comment.Text()
	}

	return strings.Trim(text, "\n")
}

func UpperSnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)

	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return strings.ToUpper(string(out))
}

func indirect(t types.Type) types.Type {
	switch t.(type) {
	case *types.Pointer:
		return indirect(t.(*types.Pointer).Elem())
	case *types.Named:
		return indirect(t.(*types.Named).Underlying())
	default:
		return t
	}
}
