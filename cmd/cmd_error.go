package cmd

import (
	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/http_error_code"
)

var (
	cmdErrorFlagErrorRegisterMethod string
)

var cmdError = &cobra.Command{
	Use:   "error",
	Short: "stringify http errors",
	Run: func(cmd *cobra.Command, args []string) {
		eg := http_error_code.NewErrorGenerator(packageName, cmdErrorFlagErrorRegisterMethod)
		eg.Output()
	},
}

func init() {
	cmdError.Flags().StringVarP(&cmdErrorFlagErrorRegisterMethod, "error-register-method", "r", "github.com/morlay/gin-swagger/http_error_code/httplib.RegisterError", "error register method")

	cmdRoot.AddCommand(cmdError)
}
