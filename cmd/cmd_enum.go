package cmd

import (
	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/swagger"
)

var cmdClientFlagRegisterEnumMethod string

var cmdEnum = &cobra.Command{
	Use:   "enum",
	Short: "stringify enum",
	Run: func(cmd *cobra.Command, args []string) {
		eg := swagger.NewEnumGenerator(packageName, cmdClientFlagRegisterEnumMethod)
		eg.Output(args...)
	},
}

func init() {
	cmdEnum.Flags().StringVarP(&cmdClientFlagRegisterEnumMethod, "register-enum-method", "r", "", "register enum method")
	cmdRoot.AddCommand(cmdEnum)
}
