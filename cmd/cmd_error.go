package cmd

import (
	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/http_error_code"
)

var (
	cmdErrorFlagErrorType string
)

var cmdError = &cobra.Command{
	Use:   "error",
	Short: "stringify http errors",
	Run: func(cmd *cobra.Command, args []string) {
		eg := http_error_code.NewErrorGenerator(packageName, cmdErrorFlagErrorType)
		eg.Output()
	},
}

func init() {
	cmdError.Flags().StringVarP(&cmdErrorFlagErrorType, "error-type", "t", "github.com/morlay/gin-swagger/http_error_code/httplib.GeneralError", "error type")

	cmdRoot.AddCommand(cmdError)
}
