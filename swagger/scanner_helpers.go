package swagger

import (
	"go/types"
	"strings"
)

func packageOfGin(packagePath string) bool {
	return strings.Contains(packagePath, "gin")
}

func typeOfGinEngine(packagePath string) bool {
	return packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "Engine"
}

func typeOfGinRouterGroup(packagePath string) bool {
	return packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "RouterGroup"
}

func typeOfGinContext(pointer *types.Pointer) bool {
	packagePath := pointer.String()
	return packageOfGin(packagePath) && getExportedNameOfPackage(packagePath) == "Context"
}

func isGinPkg(pkg *types.Package) bool {
	return pkg != nil && pkg.Name() == "gin" && packageOfGin(pkg.Path())
}

func isFuncOfGin(tpeFunc *types.Func) bool {
	return isGinPkg(tpeFunc.Pkg())
}

func isFuncWithGinContext(tpeFunc *types.Func) bool {
	if scope := tpeFunc.Scope(); scope != nil {
		for _, name := range scope.Names() {
			tpe := scope.Lookup(name).Type()
			if pointer, ok := tpe.(*types.Pointer); ok {
				return typeOfGinContext(pointer)
			}
		}
	}
	return false
}

func hasImportedGin(packages []*types.Package) bool {
	var hasGin bool

	for _, pkg := range packages {
		if isGinPkg(pkg) {
			hasGin = true
		}
	}
	return hasGin
}
