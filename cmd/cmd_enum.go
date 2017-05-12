package cmd

import (
	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/swagger"
)

var cmdEnum = &cobra.Command{
	Use:   "enum",
	Short: "stringify enum",
	Run: func(cmd *cobra.Command, args []string) {
		eg := swagger.NewEnumGenerator(packageName)
		eg.Output(args...)
	},
}

func init() {
	cmdRoot.AddCommand(cmdEnum)
}
