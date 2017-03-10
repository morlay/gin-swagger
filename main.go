package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/morlay/gin-swagger/scanner"
	"github.com/morlay/gin-swagger/swagger"
	"log"
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

	log.Printf("Generated swagger.json to %s with %d bytes\n", path, n3)
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tgin-swagger\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(aurora.Sprintf("[%s]  ", aurora.Cyan("gin-swagger")))

	flag.Usage = Usage

	sc := scanner.NewScanner(&scanner.ScannerOpts{
		PackagePath: getPackageName(),
	})

	sc.Scan()
	WriteToJSON(sc.Swagger, "swagger.json")
}
