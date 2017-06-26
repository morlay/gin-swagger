package swagger

import (
	"go/types"

	"go/ast"

	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/morlay/gin-swagger/program"
)

func NewGroup(path string, args ...ast.Expr) *Route {
	return &Route{
		Path: path,
		Args: args,
	}
}

func NewRoute(method string, path string, args ...ast.Expr) *Router {
	return &Router{
		Method: method,
		Route: Route{
			Path: path,
			Args: args,
		},
	}
}

type Route struct {
	Path string
	Args []ast.Expr
}

type Router struct {
	Route
	Method  string
	parents []*Route
}

func (router *Router) IsMethod() bool {
	return router.Method != ""
}

func (router *Router) AddGroup(r *Route) *Router {
	if router.parents == nil {
		router.parents = []*Route{}
	}
	if r != nil {
		router.parents = append(router.parents, r)
	}
	return router
}

func (router *Router) GetArgs() []ast.Expr {
	args := []ast.Expr{}

	for _, parent := range router.parents {
		if len(parent.Args) > 0 {
			args = append(parent.Args, args...)
		}
	}

	args = append(args, router.Args...)

	return args
}

func (router *Router) GetPath() string {
	prefix := ""

	for _, parent := range router.parents {
		if parent.Path != "" {
			prefix = path.Join(parent.Path, prefix)
		}
	}

	return path.Join(prefix, router.Path)
}

func (router *Router) Clone() *Router {
	return &Router{
		Method: router.Method,
		Route: Route{
			Path: router.Path,
			Args: router.Args,
		},
		parents: router.parents,
	}
}

func (router *Router) String() string {
	return router.Method + " " + router.GetPath()
}

func NewRoutesScanner() *RoutesScanner {
	return &RoutesScanner{}
}

type RoutesScanner struct {
	Routers           map[*Router]bool
	funcWithGinRouter map[*types.Func]map[*ast.CallExpr]bool
}

func (scanner *RoutesScanner) WithRouterGroup(router *Router, id *ast.Ident) {
	if id.Obj != nil {
		switch id.Obj.Decl.(type) {
		// from other function
		case *ast.Field:
			for tpeFunc, usedInCallExprMap := range scanner.funcWithGinRouter {
				if tpeFunc.Scope().Contains(id.Pos()) {
					count := 0
					baseRouter := router.Clone()

					for callExpr := range usedInCallExprMap {
						if id, ok := callExpr.Args[0].(*ast.Ident); ok {
							if count == 0 {
								scanner.WithRouterGroup(router, id)
							} else {
								scanner.WithRouterGroup(baseRouter, id)
							}
							count++
						}

					}
				}

			}
		case *ast.AssignStmt:
			assignStmt := id.Obj.Decl.(*ast.AssignStmt)
			callExpr := assignStmt.Rhs[0].(*ast.CallExpr)

			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "Group" && len(callExpr.Args) >= 1 {
					router.AddGroup(NewGroup(getRouterPathByCallExpr(callExpr), callExpr.Args[1:]...))

					if nextIdent, ok := selectorExpr.X.(*ast.Ident); ok {
						scanner.WithRouterGroup(router, nextIdent)
					}

					scanner.Routers[router] = true
				}
			}
		}
	}
}

func (scanner *RoutesScanner) collectFuncs(prog *program.Program) {
	for pkg, pkgInfo := range prog.AllPackages {
		if hasImportedGin(pkg.Imports()) {
			for _, def := range pkgInfo.Defs {
				switch def.(type) {
				case *types.Func:
					tpeFunc := def.(*types.Func)
					funcScope := tpeFunc.Scope()

					for _, name := range funcScope.Names() {
						tpeVar := funcScope.Lookup(name)
						tpeName := tpeVar.Type().String()

						if typeOfGinEngine(tpeName) || typeOfGinRouterGroup(tpeName) {
							if scanner.funcWithGinRouter == nil {
								scanner.funcWithGinRouter = map[*types.Func]map[*ast.CallExpr]bool{}
							}
							scanner.funcWithGinRouter[tpeFunc] = map[*ast.CallExpr]bool{}
						}
					}
				}
			}
		}
	}

	for tpeFunc, usedInCallExprMap := range scanner.funcWithGinRouter {
		for targetTpeFunc := range scanner.funcWithGinRouter {
			if targetTpeFunc != tpeFunc {
				for id, use := range prog.UsesInScope(targetTpeFunc.Scope()) {
					if tpeFunc.Type() == use.Type() {
						if callExpr := prog.CallExprById(id); callExpr != nil {
							if len(callExpr.Args) > 0 {
								usedInCallExprMap[callExpr] = true
							}
						}
					}
				}
			}
		}
	}
}

func (scanner *RoutesScanner) Scan(prog *program.Program) {
	scanner.collectFuncs(prog)
	scanner.Routers = map[*Router]bool{}

	for tpeFuncTarget := range scanner.funcWithGinRouter {
		funcScope := tpeFuncTarget.Scope()
		selections := prog.SelectionsInScope(funcScope)

		for selectorExpr, selection := range selections {
			if isGinMethod(selection.Obj().Name()) {
				if callExpr := prog.CallExprById(selectorExpr.Sel); callExpr != nil {
					router := NewRoute(
						selectorExpr.Sel.Name,
						getRouterPathByCallExpr(callExpr),
						callExpr.Args[1:]...,
					)

					if typeOfGinRouterGroup(selection.Recv().String()) {
						scanner.WithRouterGroup(router, selectorExpr.X.(*ast.Ident))
					}

					scanner.Routers[router] = true
				}
			}
		}
	}

	for route := range scanner.Routers {
		spew.Dump(route)
	}
}
