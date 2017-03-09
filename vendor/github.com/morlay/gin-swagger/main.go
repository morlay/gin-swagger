package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/morlay/gin-swagger/scanner"
	"github.com/morlay/gin-swagger/swagger"
)

func getPackageName() string {
	output, _ := exec.Command("go", "list").CombinedOutput()
	return strings.TrimSpace(string(output))
}

func WriteToJSON(swagger *swagger.Swagger, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		panic(err)
	}

	n3, err := f.WriteString(string(b))
	if err != nil {
		panic(err)
	}

	f.Sync()

	fmt.Printf("Generated %s with %d bytes\n", path, n3)
}

func main() {
	sc := scanner.NewScanner(&scanner.ScannerOpts{
		PackagePath: getPackageName(),
	})

	sc.Scan()
	WriteToJSON(sc.Swagger, "swagger.json")
}
