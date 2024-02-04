package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	authn "github.com/GDGVIT/configgy-cli/authn/cmd"
	"github.com/GDGVIT/configgy-cli/utils/config"
)

var cmd = &cobra.Command{
	Use:   "configgy",
	Short: "Run the configgy cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from configgy.")
		baseurl, err := cmd.Flags().GetString("baseurl")
		if err != nil {
			fmt.Println(err)
		}
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			fmt.Println(err)
		}
		if baseurl != "" {
			fmt.Println("Setting Base URL:", baseurl)
			config.SetCredentials("auth", "url", baseurl)
		}
		if token != "" {
			fmt.Println("Setting Token: ****")
			config.SetCredentials("auth", "access_token", token)
		}
	},
}

func init() {
	// add commands here with cmd.AddCommand()
	cmd.AddCommand(authn.RootCmd())
	cmd.Flags().StringP("baseurl", "b", "", "Base URL to the API")
	cmd.Flags().StringP("token", "t", "", "Token to authenticate with the API")
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
