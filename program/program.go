package program

import (
	"go/ast"
	"go/constant"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"strings"
)

type Program struct {
	*loader.Program
}

func NewProgram(packagePath string) *Program {
	ldr := loader.Config{}

	ldr.ParserMode = parser.ParseComments
	ldr.Import(packagePath)

	prog, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	return &Program{
		Program: prog,
	}
}

func (program *Program) TypeOf(e ast.Expr) types.Type {
	pkgInfo := program.Package(program.PackageOf(e).Path())

	if tpe := pkgInfo.TypeOf(e); tpe != nil {
		return tpe
	}

	return nil
}

func (program *Program) ValueOf(e ast.Expr) constant.Value {
	pkgInfo := program.Package(program.PackageOf(e).Path())

	if t, ok := pkgInfo.Types[e]; ok {
		if t.Value != nil {
			return t.Value
		}
	}
	return nil
}

func (program *Program) ScopeOf(targetNode ast.Node) *types.Scope {
	pkgInfo := program.Package(program.PackageOf(targetNode).Path())

	for _, scope := range pkgInfo.Scopes {
		if funcDecl, ok := targetNode.(*ast.FuncDecl); ok {
			if funcDecl.Body.Pos() == scope.Pos() {
				return scope
			}
		} else if targetNode.Pos() == scope.Pos() {
			return scope
		}
	}

	return nil
}

func (program *Program) ObjectOf(id *ast.Ident) types.Object {
	pkgInfo := program.Package(program.PackageOf(id).Path())
	obj := pkgInfo.ObjectOf(id)
	return obj
}

func (program *Program) DefOf(id *ast.Ident) types.Object {
	obj := program.ObjectOf(id)

	// find the defined
	switch obj.Type().(type) {
	case *types.Pointer:
		return obj.Type().(*types.Pointer).Elem().(*types.Named).Obj()
	case *types.Named:
		return obj.Type().(*types.Named).Obj()
	default:
		return obj
	}

}

func (program *Program) WhereDecl(targetTpe types.Type) ast.Expr {
	for _, pkgInfo := range program.AllPackages {
		for e, t := range pkgInfo.Types {
			if t.Type == targetTpe {
				switch e.(type) {
				case *ast.StructType:
					return e
				case *ast.Ident:
					return program.IdentOf(program.DefOf(e.(*ast.Ident)))
				case *ast.SelectorExpr:
					return program.IdentOf(program.DefOf(e.(*ast.SelectorExpr).Sel))
				}
			}
		}
	}
	return nil
}

func (program *Program) IdentOf(targetDef types.Object) *ast.Ident {
	pkgInfo := program.Package(targetDef.Pkg().Path())

	for id, def := range pkgInfo.Defs {
		if def == targetDef {
			return id
		}
	}

	return nil
}

func (program *Program) FileOf(node ast.Node) *ast.File {
	for _, pkgInfo := range program.AllPackages {
		for _, file := range pkgInfo.Files {
			if file.Pos() <= node.Pos() && file.End() > node.Pos() {
				return file
			}
		}
	}
	return nil
}

func (program *Program) PackageOf(node ast.Node) *types.Package {
	for pkg, pkgInfo := range program.AllPackages {
		for _, file := range pkgInfo.Files {
			if file.Pos() <= node.Pos() && file.End() > node.Pos() {
				return pkg
			}
		}
	}
	return nil
}

type Option struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}

func (program *Program) GetEnumOptionsByType(node ast.Node) (list []Option) {
	if ident, ok := node.(*ast.Ident); ok {
		file := program.FileOf(node)

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)

			if !ok {
				continue
			}

			if genDecl.Tok.String() == "const" {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						switch valueSpec.Type.(type) {
						case *ast.Ident:
							if valueSpec.Type.(*ast.Ident).Name == ident.Name {
								if basicLit, ok := valueSpec.Values[0].(*ast.BasicLit); ok {
									list = append(list, Option{
										Value: GetBasicLitValue(basicLit),
										Label: strings.TrimSpace(valueSpec.Comment.Text()),
									})
								}
							}
						default:
							var name = valueSpec.Names[0].Name
							if strings.HasPrefix(name, UpperSnakeCase(ident.String())) {
								var values = strings.SplitN(name, "__", 2)
								if len(values) == 2 {
									list = append(list, Option{
										Value: values[1],
										Label: strings.TrimSpace(valueSpec.Comment.Text()),
									})
								}
							}
						}

					}

				}
			}
		}
	}
	return
}

func (program *Program) CommentGroupFor(targetNode ast.Node) (commentList []*ast.CommentGroup) {
	file := program.FileOf(targetNode)

	switch targetNode.(type) {
	case *ast.StructType:
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spc := range genDecl.Specs {
					if typeSpec, ok := spc.(*ast.TypeSpec); ok {
						if typeSpec.Type == targetNode {
							if len(genDecl.Specs) > 1 {
								commentList = append(commentList, typeSpec.Doc)
							} else {
								commentList = append(commentList, genDecl.Doc)
							}
						}
					}

				}
			}
		}
	case *ast.Ident:
		ident := targetNode.(*ast.Ident)
		if ident.Obj != nil {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spc := range genDecl.Specs {
						if typeSpec, ok := spc.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == ident.Name {
								if len(genDecl.Specs) > 1 {
									commentList = append(commentList, typeSpec.Doc)
								} else {
									commentList = append(commentList, genDecl.Doc)
								}
							}
						}

					}
				}
			}
		}
	default:
		commentMap := ast.NewCommentMap(program.Fset, file, file.Comments)
		commentList = commentMap[targetNode]
	}

	return
}
