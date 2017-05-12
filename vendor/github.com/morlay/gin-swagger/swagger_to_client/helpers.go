package swagger_to_client

import (
	"fmt"
	"regexp"
	"strings"
)

func convertSwaggerPathToGinPath(str string) string {
	r := regexp.MustCompile(`/\{([^/\\}]+)\}`)
	result := r.ReplaceAllString(str, "/:$1")
	return result
}

func getRefName(str string) string {
	parts := strings.Split(str, "/")
	return parts[len(parts)-1]
}

func getPackageNameFromPath(str string) string {
	parts := strings.Split(str, ".")
	packagePaths := strings.Split(strings.Join(parts[0:len(parts)-1], "."), "vendor/")

	if len(packagePaths) > 2 {
		panic(fmt.Errorf("package name is not available `%s`", str))
	}

	if len(packagePaths) == 2 {
		return packagePaths[1]
	}

	return packagePaths[0]
}
