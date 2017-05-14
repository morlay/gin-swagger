package cmd

import (
	"github.com/spf13/cobra"

	"github.com/morlay/gin-swagger/swagger_to_client"
)

var (
	cmdClientFlagInput      string
	cmdClientFlagClientName string
	cmdClientFlagBaseClient string
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Generate go client from swagger.json",
	Run: func(cmd *cobra.Command, args []string) {
		cg := swagger_to_client.NewClientGenerator(cmdClientFlagClientName, cmdClientFlagBaseClient)
		cg.LoadSwaggerFromFile(cmdClientFlagInput)
		cg.Output()
	},
}

func init() {
	cmdClient.Flags().StringVarP(&cmdClientFlagInput, "input", "", "swagger.json", "swagger json file path")
	cmdClient.Flags().StringVarP(&cmdClientFlagClientName, "name", "", "service", "client name")
	cmdClient.Flags().StringVarP(&cmdClientFlagBaseClient, "base-client", "", "github.com/morlay/gin-swagger/swagger_to_client/client", "client name")

	cmdRoot.AddCommand(cmdClient)
}
