package cmd

import (
	"fmt"
	"os"
	"strings"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/swagger"
)

var (
	packageName string
)

func getPackageName() string {
	output, _ := exec.Command("go", "list").CombinedOutput()
	return strings.TrimSpace(string(output))
}

var cmdRoot = &cobra.Command{
	Use:   "gin-swagger",
	Short: "Generate swagger.json from gin framework codes",
	Run: func(cmd *cobra.Command, args []string) {
		sc := swagger.NewScanner(packageName)
		sc.Output("swagger.json")
	},
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cmdRoot.PersistentFlags().StringVarP(&packageName, "package", "p", getPackageName(), "package name for generating")
}
