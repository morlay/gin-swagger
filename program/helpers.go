package program

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"strconv"
	"strings"
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

func GetConstValue(value constant.Value) (uint64, error) {
	if value.Kind() != constant.Int {
		return 0, fmt.Errorf("can't happen: constant is not an integer")
	}
	i64, isInt := constant.Int64Val(value)
	u64, isUint := constant.Uint64Val(value)
	if !isInt && !isUint {
		return 0, fmt.Errorf("internal error: value of %s is not an integer", value.String())
	}
	if !isInt {
		u64 = uint64(i64)
	}
	return u64, nil
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

func IsVendorPackage(path string) bool {
	return len(strings.Split(path, "vendor/")) > 1
}

func IsSubPackageOf(basePath string, path string) bool {
	return strings.Index(path, basePath) == 0
}
