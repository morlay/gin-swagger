package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/morlay/gin-swagger/client_generator"
	"github.com/morlay/gin-swagger/scanner"
)

func getPackageName() string {
	output, _ := exec.Command("go", "list").CombinedOutput()
	return strings.TrimSpace(string(output))
}

func main() {
	log.SetPrefix(aurora.Sprintf("[%s]  ", aurora.Cyan("gin-swagger")))

	flags := flag.NewFlagSet("gin-swagger", flag.ContinueOnError)

	err := flags.Parse(os.Args[1:])

	if err != nil {
		panic(err)
	}

	switch flags.Arg(0) {
	case "client":
		input := flag.String("input", "swagger.json", "swagger json file path")
		clientName := flag.String("name", "service", "client name")
		baseClient := flag.String("base", "github.com/morlay/gin-swagger/example/client", "client name")

		cg := client_generator.NewClientGenerator(*clientName, *baseClient)
		cg.LoadSwaggerFromFile(*input)
		cg.Output()
	default:
		sc := scanner.NewScanner(&scanner.ScannerOpts{
			PackagePath: getPackageName(),
		})
		sc.Output("swagger.json")
	}
}
